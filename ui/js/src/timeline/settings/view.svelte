<script>
import { settings,state } from './store.js'
import { fade} from 'svelte/transition'
import {onMount, onDestroy, createEventDispatcher} from 'svelte'

import Info from './tabs/info.svelte'
import Appearance from './tabs/appearance.svelte'
import Pages from './tabs/pages/pages.svelte'

const dispatch = createEventDispatcher()

let active = false;

onMount(() => {
  fetchState().then((res) => {
    console.log(res)
    if(res?.state) {
      active = true
    }
    $state = res.state
  }).then(() => {
    let title = $state?.filter(x => x.type == 'm.room.name')[0]?.content["name"]
    let about = $state?.filter(x => x.type == 'm.room.topic')[0]?.content["topic"]
    let avatar = $state?.filter(x => x.type == 'm.room.avatar')[0]?.content["url"]
    let css = $state?.filter(x => x.type == 'com.hummingbard.room.style')[0]?.content["css"]
    let header = $state?.filter(x => x.type == 'com.hummingbard.room.header')[0]?.content["url"]
    $settings = {
      info: {
        title: title,
        about: about,
        avatar: avatar,
      },
      appearance: {
        header: header,
        css: css,
      }
    }
    console.log("store settings updated to", $settings)
  })
})



async function fetchState() {
    let endpoint = `/room/state`

    let data = {
        room_id: window.timeline.room_id,
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
    if (resp.ok) { // if HTTP-status is 200-299
    } else {
      alert("HTTP-Error: " + resp.status);
    }
    const ret = await resp.json()
    return Promise.resolve(ret)
}

function kill() {
  dispatch('kill', true)
}







async function saveState() {
    let endpoint = `/room/info/update`

    let data = {
      room_id: window.timeline.room_id,
      profile: window.timeline.profile,
      info: {},
      appearance: {},
    }

    let title = $state?.filter(x => x.type == 'm.room.name')[0]?.content["name"]
    let about = $state?.filter(x => x.type == 'm.room.topic')[0]?.content["topic"]
    let avatar = $state?.filter(x => x.type == 'm.room.avatar')[0]?.content["url"]
    let css = $state?.filter(x => x.type == 'com.hummingbard.room.style')[0]?.content["css"]
    let header = $state?.filter(x => x.type == 'com.hummingbard.room.header')[0]?.content["url"]

    data.info.title = $settings.info.title

    data.info.about = $settings.info.about

    data.info.avatar = $settings.info.avatar

    data.appearance.header = $settings.appearance.header

    data.appearance.css = $settings.appearance.css

  if(!data.info && !data.appearance) {
    console.log("nothing to upadte")
    return
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
    if (resp.ok) { // if HTTP-status is 200-299
    } else {
      alert("HTTP-Error: " + resp.status);
    }
    const ret = await resp.json()
    return Promise.resolve(ret)
}

let saving = false;

function save() {
  saving = true
  saveState().then((res) => {
    console.log(res)
    if(res?.updated) {
      saving = false
      kill()
    }
  }).then(() => {
    location.reload()
  })
}

let tab = 'info';

$: infoTab = tab == 'info'
$: appTab = tab == 'appearance'
$: perTab = tab == 'permissions'
$: subTab = tab == 'sub-spaces'
$: pagesTab = tab == 'pages'

function switchTab(e) {
  tab = e.target.id
}

function activeTab(e) {
  return tab == e.target.id
}



$: profile = window.timeline?.profile

let SubSpaces;
let subSpacesLoaded;
$: if(subTab) {
  import('./tabs/spaces/sub-spaces.svelte').then(res => {
    SubSpaces = res.default
    subSpacesLoaded = true
  })
}

</script>



{#if active}
<div class="modal-container main ph3" transition:fade="{{duration: 73}}">

  <div class="modal-inner-np start flex flex-column " >

    <div class="pa3 flex flex-column">

      <div class="flex no-select">
        <div id="info" class="tab-item" class:active-tab={infoTab}
          on:click={switchTab}>Info</div>
        <div id="appearance" class="tab-item" class:active-tab={appTab}
          on:click={switchTab}>Appearance</div>
        <div id="sub-spaces" class="tab-item" class:active-tab={subTab}
          on:click={switchTab}>Sub-Spaces</div>
        {#if !profile}
        <div id="pages" class="tab-item" class:active-tab={pagesTab}
          on:click={switchTab}>Pages</div>
        {/if}
      </div>

      <div class="tab-view flex mt4">

      {#if tab == 'info'}
        <Info />
      {/if}

      {#if tab=='appearance'}
        <Appearance />
      {/if}


      {#if tab=='sub-spaces' && subSpacesLoaded}
        <SubSpaces/>
      {/if}


      {#if tab=='pages' && !profile}
        <Pages/>
      {/if}


      </div>


        <div class="flex mt4">
          <div class="">
              <button on:click={save} disabled={saving}>Save</button>
              {#if !saving}
                <button class="light" on:click={kill}>Cancel</button>
                {/if}
          </div>
            {#if saving}
                <div class="lds-ring gr-center ml3"><div></div><div></div><div></div><div></div></div>
            {/if}


        </div>




    </div>


  </div>

  <div class="mask absolute" on:click={kill}></div>

</div>
{/if}


<style>
.modal-container {
    top: 0;
    left: 0;
    position: fixed;
    width: 100%;
    height: 100%;
    z-index: 49999;
    display: grid;
    grid-template-columns: auto;
    grid-template-rows: auto;
}

.modal-inner-np {
    justify-self: center;
    display: grid;
    grid-template-columns: auto;
    grid-template-rows: auto;
    z-index: 5000;

    background-color: var(--m-bg);
    -webkit-box-shadow: 0px 4px 24px -1px rgba(0,0,0,0.05);
    -moz-box-shadow: 0px 4px 24px -1px rgba(0,0,0,0.05);
    box-shadow: 0px 4px 24px -1px rgba(0,0,0,0.05);
    border-radius: 17px;
    transition: 0.1s;
    word-break: break-word;
}

.start {
    align-self: center;
    width: 680px;
}

.tab-item {
  background: var(--m-bg);
  color: var(--text);
  padding: 0;
  padding-bottom: 0.5rem;
  margin: 0;
  margin-right: 1rem;
  font-size: 0.9rem;
  cursor: pointer;
  border-bottom: 1px transparent;
}

.tab-item:hover {
  border-bottom: 1px solid var(--primary-gray);
}

.active-tab {
  border-bottom: 1px solid var(--primary-dark);
  font-weight: bold;
}

.tab-view {
  min-height: 60vh;
}

@media screen and (max-width: 680px) {
  .start {
    width: 100%;
  }
}

.mask {
    top: 0;
    left: 0;
    height: 100%;
    width: 100%;
    background: var(--mask)
}



</style>
