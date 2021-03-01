<script>
export let type;
export let id;
export let name = '';
export let alias;
export let inline;

let joining = false;

function isMember() {
    let ind = identity?.joined_rooms.findIndex(x => x.room_id == id)
    if(ind != -1) {
        return true
    }
    return false
}

$: member = inline ? isMember() : window.timeline?.member

$: joinText = (type == `user`) ? `Follow ${name}` : `Join ${name}`

$: joinText = (type == `user`) ? member ? `Unfollow ${name}` : `Follow ${name}` : member ? `Leave ${name}` : `Join ${name}`



async function joinRoom() {
    let endpoint = `/room/join`

    if(member) {
        endpoint = `/room/leave`
    }

    let data = {
        id: id,
        alias: alias,
    }

    let options = {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Content-Type': 'application/json'
        }
    }


    if(authenticated && identity?.access_token) {
        options.headers['Authorization'] = identity.access_token
    }
    console.log(options)

    let resp = await fetch(endpoint, options)
    const ret = await resp.json()
    return Promise.resolve(ret)
}

function join() {
    joining = true
  joinRoom().then((res) => {
    console.log("joined?: ",res)
      if(res?.joined) {
          if(!inline) {
              location.reload()
          } else {
              member = true
          }
      } else if(res?.left) {
          if(!inline) {
              location.reload()
          } else {
              member = false
          }
      }
  }).then(() => {
      joining = false
  })
}

</script>
<button class="" class:light={member} on:click={join} disabled={joining}>{joinText}</button>
