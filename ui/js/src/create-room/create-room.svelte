<script>
import {fade} from 'svelte/transition'
import {onMount,createEventDispatcher} from 'svelte'
import {debounce} from '../utils/utils.js'

const dispatch = createEventDispatcher();

export let sub = false;
export let level;
export let roomID;

let active = false;

let usernameInput;
let titleInput;
let aboutInput;
let typeInput;

let username = '';

function resetUsername() {
  usernameInput.value = ''
  username = ''
  available = false;
  usernameAvailable = true;
  usernameInput.focus()
}

let privateRoom = false;
let nsfw = false;


let checking = false;
let available = false;
let usernameAvailable = true;

onMount(() => {
  active = true
})

$: if(active && usernameInput) {
  usernameInput.focus()
}

let locked = false;

function create() {
  if(!usernameAvailable || !available) {
    usernameInput.focus()
    return
  }
  if(usernameInput.value.length == 0) {
    usernameInput.focus()
    return
  }
  if(titleInput.value.length == 0) {
    titleInput.focus()
    return
  }
  if(aboutInput.value.length == 0) {
    aboutInput.focus()
    return
  }

  locked = true
    createRoom().then((res) => {
      console.log(res)
      if(res?.created) {
        if(!sub) {
          if(identity.federated) {
            let path = res?.room?.canonical_alias?.substring(1)
            if(path) {
              window.location.href = `/${path}`
            }
          } else {
            window.location.href = `/${res?.room?.alias}`
          }
        } else {
          dispatch('created', res?.room)
        }
      }
    }).then(() => {
    })
}

function updateUsername(e) {
  const letters = /^[0-9a-zA-Z-]+$/;
  if(!e.key.match(letters)){
    e.preventDefault()
  }
  usernameAvailable = true
  available = false
  if(username.length === 0) {
      checking = false
      return
  }
}

function reset() {
  available = false
  debounce(() =>{
    checking = true
    if(usernameInput.value.length == 0) {
      checking = false
      usernameAvailable = true
      available = false
      return
    }
    checkUsername().then((res) => {
      console.log(res)
        if(res?.available) {
          usernameAvailable = true
          available = true
        } else if(!res?.available) {
          usernameAvailable = false
          available = false
        }
        checking = false
    }).then(() => {
    })

  }, 500, this)
}

$: domain = sub ? `https://${identity.home_server}/${level}` : `https://${identity.home_server}`

$: usernameText = username.length > 0 ? username : `music`


async function checkUsername() {

    let endpoint = `/username/available`

    let data = {
        username: username,
    };

  if(sub) {
    data.username= username
    data.parent_room_id = roomID
    data.sub_space = true
  }

    let resp = await fetch(endpoint, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Authorization': identity.access_token,
            'Content-Type': 'application/json'
        }
    })

    const ret = await resp.json()
    return Promise.resolve(ret)

}

async function createRoom() {

    let endpoint = `/create`

    let data = {
      username: username,
      title: titleInput.value,
      about: aboutInput.value,
      type: typeInput.value,
      nsfw: nsfw,
    };

  if(sub) {
    data.parent_room_id = roomID
    data.sub_space = true
  }

  console.log(data)

    let resp = await fetch(endpoint, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Authorization': identity.access_token,
            'Content-Type': 'application/json'
        }
    })

    const ret = await resp.json()
    return Promise.resolve(ret)

}


function cancel() {
  dispatch('cancel', true)
}

</script>

{#if active}
  <div class="gr-center cr brd" 
  class:w-100={sub} 
  class:sub={sub} 
  transition:fade="{{duration: 33}}">

  <div class="flex flex-column lh-copy" 
     class:pa3={!sub}>
    {#if !sub}
    <div class="mb3">
      <span class="tab">
        <strong>Create</strong>
      </span>
    </div>
    {/if}

    {#if sub}
    <div class="mb3">
      <span class="tab">
        <strong>Parent: {level}</strong>
      </span>
    </div>
    {/if}

    <div class="flex flex-column" class:mt3={!sub}>
      <div class="">
        <span class="small bold">Name</span>
      </div>
      <div class="mt2 relative">
        <input
          class:oops={!usernameAvailable}
          placeholder="music"
          on:keypress={updateUsername}
          on:input={reset}
          bind:this={usernameInput}
          bind:value={username}
        />
        {#if checking}
          <div class="checking mh2 gr-default">
            <div class="lds-ring gr-center"><div></div><div></div><div></div><div></div></div>
          </div>
        {/if}
        {#if !usernameAvailable}
          <div class="checking mh2 gr-default pointer" on:click={resetUsername}>
            <svg class="gr-center" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M5.72 5.72a.75.75 0 011.06 0L12 10.94l5.22-5.22a.75.75 0 111.06 1.06L13.06 12l5.22 5.22a.75.75 0 11-1.06 1.06L12 13.06l-5.22 5.22a.75.75 0 01-1.06-1.06L10.94 12 5.72 6.78a.75.75 0 010-1.06z"></path></svg>
          </div>
        {/if}
        {#if available}
          <div class="checking mh2 gr-default">
            <svg class="gr-center" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M21.03 5.72a.75.75 0 010 1.06l-11.5 11.5a.75.75 0 01-1.072-.012l-5.5-5.75a.75.75 0 111.084-1.036l4.97 5.195L19.97 5.72a.75.75 0 011.06 0z"></path></svg>
          </div>
        {/if}
      </div>
      <div class="mt2 fs-09">
      {#if usernameAvailable}
        <span class="o-80">{domain}</span>/<span class="">{usernameText}</span>
      {:else}
        <span class="primary">That name is not available</span>
      {/if}
      </div>
    </div>

    <div class="mt3 flex flex-column">
      <div class="">
        <span class="small bold">Title</span>
      </div>
      <div class="mt2">
        <input
          placeholder="Great Music"
          bind:this={titleInput}
        />
      </div>
    </div>

    <div class="mt3 flex flex-column">
      <div class="">
        <span class="small bold">About</span>
      </div>
      <div class="mt2">
        <textarea
          style="height:100px;"
          placeholder="A place to share your favorite music."
          bind:this={aboutInput}
        ></textarea>
      </div>
    </div>

    <div class="mt3 flex">

      <div class="fs-09">
        <select class="sel" name="type" bind:this={typeInput}>
            <option value="space">space</option>
            <option value="gallery">gallery</option>
        </select>
      </div>


      <div class="flex gr-center ml3 no-select">
          <input id="nsfw" type=checkbox bind:checked={nsfw} hidden>
          <label class="label" class:checked-nsfw={nsfw} for="nsfw">
            {nsfw ? 'NSFW' : 'SFW'}
          </label>
      </div>

      <div class="flex-one"></div>

      {#if sub}
      <div class="gr-center mr3">
        <button class="light" on:click={cancel}>Cancel</button>
      </div>
      {/if}

      {#if locked}
      <div class="gr-center mr3">
        <div class="lds-ring"><div></div><div></div><div></div><div></div></div>
      </div>
      {/if}

      <div class="gr-center">
        <button on:click={create} disabled={locked}>Create</button>
      </div>
    </div>

  </div>

</div>
{/if}

<style>

.sub {
  padding: 1rem;
  border-radius: 17px;
  border: 1px solid var(--primary-gray);
}

.label {
  font-size: small;
  padding: 0.125rem 0.25rem;
  cursor: pointer;
  background-color: var(--primary-light-gray);
  border-radius: 2px;
  transition: 0.1s;
}

.label:hover {
  font-weight: bold;
}

.checked {
  font-weight: bold;
  background-color: var(--primary-dark);
  color: var(--white);
}

.checked-nsfw {
  font-weight: bold;
  background-color: var(--primary);
  color: var(--white);
}

.cr {
  min-width: 500px;
}

.tab {
  border-radius: 500px;
  background-color: var(--yellow);
  padding: 0.25rem 0.5rem;
  font-size: 0.8rem;
}

.oops {
  border: 1px solid red;
}

.checking {
  position: absolute;
  top: 0;
  right: 0;
  height: 100%;
}

  .sel {
    border-radius: 4px;
  }

input, textarea {
  width: 100%;
}

@media screen and (max-width: 538px) {
  .cr {
    min-width: 100%;
  }
}
</style>
