<template>
  <div class="container">
    <add-user :ws="ws"></add-user>
    <br>
    <all-users></all-users>
  </div>
</template>

<script>

import AddUser from './AddUser'
import AllUsers from './AllUsers'
var ws

export default {
  name: 'Home',
  data() {
    return {
      ws
    }
  },
  created() {
    this.ws = new WebSocket("ws://localhost:1323/users")

    this.ws.onopen = function() {
      console.log('Connected')
    }
  },
  components: {
    AddUser,
    AllUsers,
  }
}

</script>
