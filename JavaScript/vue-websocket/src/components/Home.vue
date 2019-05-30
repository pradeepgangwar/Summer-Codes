<template>
  <div class="container">
    <add-user></add-user>
    <br>
    <all-users></all-users>
  </div>
</template>

<script>

import AddUser from './AddUser'
import AllUsers from './AllUsers'

export default {
  name: 'Home',
  components: {
    AddUser,
    AllUsers,
  }
}

</script>
