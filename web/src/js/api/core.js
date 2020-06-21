import axios from "axios"

const api = {
  getSession: function() {
    return axios({
      method: "get",
      url: "/api/session",
    })
  },

  newSession: function(numOfPlayers) {
    return axios({
      method: "post",
      url: `/api/new_game?numOfPlayers=${numOfPlayers}`,
    })
  },

  playTurn: function(index) {
    return axios({
      method: "post",
      url: `/api/play_turn/${index}`
    })
  }
  
}

export default api
