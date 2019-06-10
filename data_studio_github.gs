var connector = connector || {};

function getAuthType() {
  var response = {type: 'NONE'};
  return response;
}

function getConfig() {
  var cc = DataStudioApp.createCommunityConnector();
  var config = cc.getConfig();

  config
    .newInfo()
    .setId('instructions')
    .setText(
      'Enter github username to fetch their information.'
    );

  config
    .newTextInput()
    .setId('username')
    .setName(
      'Enter a single username.'
    )
    .setAllowOverride(true);

  return config.build();
}

function getFields() {
  var cc = DataStudioApp.createCommunityConnector();
  var fields = cc.getFields();
  var types = cc.FieldType;

  fields
    .newDimension()
    .setId('username')
    .setName('Username')
    .setType(types.TEXT);

  fields
    .newDimension()
    .setId('type')
    .setName('Type')
    .setType(types.TEXT);

  fields
    .newMetric()
    .setId('count')
    .setName('Count')
    .setType(types.NUMBER)

  return fields;
}

function getSchema(request) {
  var fields = getFields().build();
  return {schema: fields};
}

function getData(request) {

  var requestedFieldIds = request.fields.map(function(field) {
    return field.name;
  });
  var requestedFields = getFields().forIds(requestedFieldIds);
  try {
    var apiResponse = connector.fetchDataFromApi(request);
  } catch (e) {
    connector.throwError('Unable to fetch data from source.', true);
  }
  var validResponse = JSON.parse(apiResponse);

  try {
    var data = connector.getFormattedData(validResponse, requestedFields);
  } catch (e) {
    connector.throwError('Unable to process data in required format.', true);
  }
  console.log(data);

  return {
    schema: requestedFields.build(),
    rows: data
  };
}

function isAdminUser() {
  return true;
}

connector.getFormattedData = function(parsedResponse, requestedFields) {
  var types = ['followers', 'following', 'repos'];
  var data = [];
  var formatted_data = types.map(function (type) {
    return connector.formatData(
      type,
      parsedResponse
    );
  });
  data = data.concat(formatted_data)
  return data;
};

connector.fetchDataFromApi = function(request) {
  var url = [
    'https://api.github.com/users/',
    request.configParams.username
  ];
  var response = UrlFetchApp.fetch(url.join(''));
  return response;
};

connector.formatData = function(type, response) {
  var values = [];
  values.push(type);
  switch (type) {
    case 'followers':
      values.push(response.followers);
      break;
    case 'following':
      values.push(response.following);
      break;
    case 'repos':
      values.push(response.public_repos);
      break;
    default:
      values.push('');
  }
  return {values: values};
};


connector.throwError = function(message, userSafe) {
  if (userSafe) {
    message = 'DS_USER:' + message;
  }
  throw new Error(message);
};
