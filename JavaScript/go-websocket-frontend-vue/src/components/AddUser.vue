  
<template>
  <div class="row">
    <div class="col-sm-4 offset-sm-4 text-center">
      <b-form @submit.prevent="adduser">

        <b-form-group id="name"
                      label="Your Name:"
                      label-for="nameinput">
            <b-form-input id="nameinput"
                        type="text"
                        v-model="name"
                        required
                        placeholder="Enter Name">
            </b-form-input>
        </b-form-group>
        <b-form-group id="email"
                      label="Email:"
                      label-for="emailinput">
            <b-form-input id="emailinput"
                        type="text"
                        v-model="email"
                        required
                        placeholder="Enter Email">
            </b-form-input>
        </b-form-group>
            <b-form-group id="about"
                      label="About:"
                      label-for="aboutinput">
            <b-form-input id="aboutinput"
                        type="text"
                        v-model="about"
                        required
                        placeholder="About">
            </b-form-input>
        </b-form-group>

        <b-button type="submit" variant="primary">Add User</b-button>
        <br>
      </b-form>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'AddUser',
  props:['ws'],
  data () {
    return {
      name: '',
      email: '',
      about: ''
    }
  },
  methods: {
    adduser () {
        var body = JSON.stringify({
            name: this.$data.name,
            email: this.$data.email,
            about: this.$data.about
        })
        this.ws.send(body)
    }
  }
}
</script>