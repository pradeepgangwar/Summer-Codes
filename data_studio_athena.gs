var cc = DataStudioApp.createCommunityConnector();

var userToken = 'KEY';
var apiUrl = 'https://deloitte-api.atlan.com';
var pingUrl = 'https://deloitte-api.atlan.com/organizations/current';
var tableApi = 'https://deloitte-api.atlan.com/v1/statement';

function isAdminUser() {
  return false;
}

function getAuthType() {
  var authTypes = cc.AuthType;
  return cc
    .newAuthTypeResponse()
    .setAuthType(authTypes.KEY)
    .build();
}

function setCredentials(request) {
  var token = request.key;
  
  var credentials = validateCredentials(token);
  if (credentials === false) {
    return {
      errorCode: 'INVALID_CREDENTIALS'
    };
  }
  var userProperties = PropertiesService.getUserProperties();
  userProperties.setProperty(userToken, token);
  return {
    errorCode: 'NONE'
  };
}

function isAuthValid() {
  var atlanToken = getToken();
  return validateCredentials(atlanToken);
}

function getConfig() {
  var config = cc.getConfig();

  config
    .newInfo()
    .setId('instructions')
    .setText(
      'Enter the table name from athena to fetch the data and visualize.'
    );

  config
    .newTextInput()
    .setId('tablename')
    .setName(
      'Enter tablename'
    )

  return config.build();
}

function getSchema(request) {
  var response  = fetchTableUrl(request);
  if (response === false) {
    throwConnectorError('Error fetchinh data from the athena table', true);
  }
  var nextUrl = response.nextUri;
  nextUrl = nextUrl.substring(nextUrl.lastIndexOf(':') + 5);
  nextUrl = apiUrl + nextUrl;
  var token  = getToken();
  var options = {
    headers: {
      APIKEY: token
    }
  };
  if (nextUrl !== undefined) {
    while(true) {
      response = getApiResponse(nextUrl, options);
      if (response.hasOwnProperty('columns')) {
        var columns = response.columns;
        break;
      }
      if (!response.hasOwnProperty('nextUri')) {
        break;
      }
      nextUrl = response.nextUri;
      nextUrl = nextUrl.substring(nextUrl.lastIndexOf(':') + 5);
      nextUrl = apiUrl.concat(nextUrl);
    }
  }
  // Cache Key for schema
  var schemaCacheKey = nextUrl.concat('-schema');
  var cache = CacheService.getScriptCache();
  var cachedSchema = JSON.parse(cache.get(schemaCacheKey));
  if (cachedSchema === null) {
    var schema = buildSchema(columns, schemaCacheKey);
  } else {
    var schema = cachedSchema;
  }
  return {
    schema: schema
  };  
}

function getData(request) {
  var cache = CacheService.getScriptCache();
  var response = fetchTableUrl(request);
  if (response === false) {
    throwConnectorError('Error fetching data from the athena table', true);
  }
  // Particular vars
  var tableRows = [];
  var schemaCacheKey = null;
  var schema = null;
  var columns = null;
  
  var nextUrl = response.nextUri;
  nextUrl = nextUrl.substring(nextUrl.lastIndexOf(':') + 5);
  nextUrl = apiUrl + nextUrl;
  var token  = getToken();
  var options = {
    headers: {
      APIKEY: token
    }
  };
  if (nextUrl !== undefined) {
    while(true) {
      response = getApiResponse(nextUrl, options);
      // Check for cached schema
      if (schemaCacheKey === null && response.hasOwnProperty('columns')) {
        schemaCacheKey = nextUrl.concat('-schema');
        var cachedSchema = JSON.parse(cache.get(schemaCacheKey));
        if (cachedSchema === null) {
          schema = buildSchema(response.columns);
          columns = response.columns;
        } else {
          schema = cachedSchema;
          columns = response.columns;
        }
      }
      if (response.hasOwnProperty('data')) {
        var rows = response.data;
        rows.map(function (row) {
          tableRows.push(row);
        });
      }
      if (!response.hasOwnProperty('nextUri')) {
        break;
      }
      nextUrl = response.nextUri;
      nextUrl = nextUrl.substring(nextUrl.lastIndexOf(':') + 5);
      nextUrl = apiUrl.concat(nextUrl);
    }
  }
  var requestedSchema = request.fields.map(function(field) {
    for (var i = 0; i < schema.length; i++) {
      if (schema[i].name === field.name) {
        return schema[i];
      }
    }
  });
//  console.log(requestedSchema);
  var requestedData = processData(tableRows, columns, requestedSchema);
  return {
    schema: requestedSchema,
    rows: requestedData
  };
}

function getApiResponse(url, options) {
  // Initialize cache service
  var cache = CacheService.getScriptCache();
  var cachedResponse = JSON.parse(cache.get(url));
  if (cachedResponse === null) {
    try {
      var tableResponse = UrlFetchApp.fetch(url, options);
      cache.put(url, tableResponse);
    } catch (err) {
      console.log(err);
    }
  } else {
    var tableResponse = cachedResponse;
  }
  tableResponse = JSON.parse(tableResponse);
  return tableResponse;
}

function processData(data, columns, schema) {
  var dataIndexes = schema.map(function(field) {
    for (var i=0; i<columns.length; i++) {
      var name = (columns[i].name).replace(/\s+/g, '');
      if (name === field.name) {
        return Math.floor(i);
      }
    }
  });
  console.log(dataIndexes);
  var result = [];
  for (var rowIndex = 0; rowIndex < data.length; rowIndex++) {
    var row = data[rowIndex];
    var rowData = dataIndexes.map(function(columnIndex) {
      return row[columnIndex];
    });
    result.push({
      values: rowData
    });
  }
  return result;
}

function buildSchema(columns, cacheKey) {
  var schema = [];
  for (var i=0; i<columns.length; i++) {
    var fieldSchema = mapColumn(columns[i]);
    schema.push(fieldSchema);
  }
  var cache = CacheService.getScriptCache();
  cache.put(cacheKey, JSON.stringify(schema));
  return schema;
}

function mapColumn(column) {
  var field = {};
  var name = column.name;
  name = name.replace(/\s+/g, '');
  field.name = name;
  field.label =name;
  if (column.type === 'varchar' || column.type === 'varchar(32)') {
    field.dataType = 'STRING';
    field.semantics = {};
    field.semantics.conceptType = 'DIMENSION';
  } else if (column.type === 'timestamp') {
    field.dataType = 'STRING';
    field.semantics = {};
    field.semantics.conceptType = 'DIMENSION';
    field.semantics.semanticType = 'YEAR_MONTH_DAY';
  } else {
    field.dataType = 'NUMBER';
    field.semantics = {};
    field.semantics.conceptType = 'METRIC';
  }
  return field;
}

function fetchTableUrl(request) {
  var token  = getToken();
  var data = "SELECT * from " + request.configParams.tablename;
  var options = {
    'method' : 'POST',
    'headers': {
      APIKEY: token
    },
    'contentType': 'application/json',
    'payload' : data
  };
  try {
    var response = UrlFetchApp.fetch(tableApi, options);
  } catch (err) {
    return false;
  }
  if (response.getResponseCode() == 200) {
    return JSON.parse(response);
  }
  return false;
}

function validateCredentials(token) {
  if (token === '') {
    return false;
  }
  var options = {
    headers: {
      APIKEY: token
    }
  };
  try {
    var response  = UrlFetchApp.fetch(pingUrl, options);
  } catch (err) {
    return false;
  }
  
  if (response.getResponseCode() == 200) {
    return true;
  }
  return false;
}

function getToken() {
  var properties = PropertiesService.getUserProperties();
  var token = properties.getProperty(userToken);
  return token;
}

function resetAuth() {
  var userProperties = PropertiesService.getUserProperties();
  userProperties.deleteProperty(userToken);
}


function throwConnectorError(message, userSafe) {
  userSafe =
    typeof userSafe !== 'undefined' && typeof userSafe === 'boolean'
      ? userSafe
      : false;
  if (userSafe) {
    message = 'DS_USER:' + message;
  }
  throw new Error(message);
}
