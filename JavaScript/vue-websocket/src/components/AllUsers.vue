<template>
  <div class="row">
    <div class="col-sm-10 offset-sm-1">
      <table class="table table-hover table-success" v-if="users.length > 0">
        <thead>
          <tr>
            <th scope="col" style="text-align:center;">Name</th>
            <th scope="col" style="text-align:center;">Email</th>
            <th scope="col" style="text-align:center;">About</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="user in users"
            :key="user.id">
            <td align="center">{{ user.name }}</td>
            <td align="center">{{ user.email }}</td>
            <td align="center">{{ user.about }}</td>            
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
    data() {
        return {
            users: [],
            loading: true
        }
    },

    mounted () {
    axios
        .get('http://localhost:1323/users')
        .then(users => {
            console.log(users.data)
            this.users = users.data
        })
        .catch(error => {
        console.log(error)
        })
    }
}
</script>
