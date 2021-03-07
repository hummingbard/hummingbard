<script>
import {tick, createEventDispatcher} from 'svelte'
const dispatch = createEventDispatcher();
import {debounce} from '../utils/utils.js'

import {addCollection} from './store.js'
import { makeid } from '../utils/utils.js'

let nameInput;
let name = ''
let descriptionInput;
let description = ''

function updateDetails(e) {
}

function kill() {
  dispatch('kill')
}

async function save() {
  if(name.length == 0) {
    alert("Name cannot be empty.")
    await tick();
    nameInput.focus()
    return
  }
  if(description.length == 0) {
    alert("Description cannot be empty.")
    await tick();
    descriptionInput.focus()
    return
  }

  saveCollection().then((res) => {
    console.log(res)
    if(res) {
      addCollection({
        id: makeid(16),
        name: name,
        path: slugged,
        description: description,
      })
      kill()
    }
  }).then(() => {
  })

}

async function saveCollection() {
    let endpoint = `/collections/create`
    let data = {
      name: name,
      path: slugged,
      description: description,
    };
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

let checking = false;
let available = false;
let nameAvailable = true;

function updateName(e) {
  const letters = /^[0-9a-zA-Z- ]+$/;
  if(!e.key.match(letters)){
    e.preventDefault()
  }
  nameAvailable = true
  available = false
  if(name.length === 0) {
      checking = false
      return
  }
}

function reset() {
  available = false
  debounce(() =>{
    checking = true
    if(nameInput.value.length == 0) {
      checking = false
      nameAvailable = true
      available = false
      return
    }
    checkName().then((res) => {
      console.log(res)
        if(res?.available) {
          nameAvailable = true
          available = true
        } else if(!res?.available) {
          nameAvailable = false
          available = false
        }
        checking = false
    }).then(() => {
    })

  }, 500, this)
}


$: slugged = name.replaceAll(" ", "")

async function checkName() {
    let endpoint = `/collections/available`
    let data = {
        name: slugged,
    };
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

function resetName() {
  nameInput.value = ''
  name = ''
  available = false;
  nameAvailable = true;
  nameInput.focus()
}

$: domain =  `https://${identity.home_server}`
$: nameText = name.length > 0 ? slugged : `myspacecollection`

</script>

<div class="gr-center flex flex-column w-100 pa3 brd-btm">

  <div class="">
    <span class=""><strong>New Collection<strong></span>
  </div>

  <div class="mt4">
    <div class="">
      <span class="small"><strong>Name</strong></span>
    </div>
    <div class="mt2 fs-09 relative">
      <input placeholder="My Space Collection" 
          class:oops={!nameAvailable}
        on:keypress={updateName}
        on:input={reset}
        bind:this={nameInput}
        bind:value={name}
      autofocus/>
        {#if checking}
          <div class="checking mh2 gr-default">
            <div class="lds-ring gr-center"><div></div><div></div><div></div><div></div></div>
          </div>
        {/if}
        {#if !nameAvailable}
          <div class="checking mh2 gr-default pointer" on:click={resetName}>
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
      {#if nameAvailable}
        <span class="o-80">{domain}</span>/<span class="">{nameText}</span>
      {:else}
        <span class="primary">That name is not available</span>
      {/if}
      </div>
  </div>

  <div class="mt3">
    <div class="">
      <span class="small"><strong>Description</strong></span>
    </div>
    <div class="mt2 fs-09">
      <textarea 
        bind:this={descriptionInput}
        style="height:120px;"
        bind:value={description}
        on:keypress={updateDetails}
        placeholder="Description"
        max-length="500"
        />
    </div>
  </div>


  <div class="mt4 flex">
    <div class="flex-one"></div>

    <div class="">
      <button class="" on:click={save}>Save</button>
      <button class="light" on:click={kill}>Cancel</button>
    </div>
  </div>

</div>


<style>
.oops {
  border: 1px solid red;
}

.checking {
  position: absolute;
  top: 0;
  right: 0;
  height: 100%;
}

</style>
