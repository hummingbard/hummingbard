<script>
import {fade} from 'svelte/transition'
import tippy from 'tippy.js'
import {onMount} from 'svelte'
import { getPostById } from '../../timeline/store.js'

export let id;

let icon;
let content;
let menu;

$: post = getPostById(id)

onMount(() => {
      load()
})


function load() {
    menu = tippy(icon, {
        content: content,
        theme: 'menu',
        placement: "bottom",
        duration: "40",
        animation: "shift-away" ,
        trigger: "click",
        arrow: true,
        interactive: true,
        onHide(menu) {
        },
        onShow(menu) {
        },
    });

    content.classList.remove('dis-no')
}

$: admin = (authenticated && window.timeline?.admin) || identity?.user_id == '@hummingbard.com:hummingbard.com'

$: owner = authenticated && window.timeline?.owner 

$: sender = authenticated && identity?.user_id == post?.sender

$: showRedact = admin || owner || sender


function redact() {
  let el = document.querySelector(`#${post.short_id}`)
  if(el) {
    const oh = el.offsetHeight
    el.style.height = `${oh}px`
    el.style.transition = `0.2s`
    el.innerHTML = `
        <div class="gr-default h-100 w-100">
            <div class="gr-center tc fs-09 bold">
                Post Deleted.
            </div>
        </div>
    `
    setTimeout(() => {
        el.innerHTML = ``
        el.style.height = `0px`
    }, 1000)
    setTimeout(() => {
        el.remove()
    }, 4000)
  }
  redactPost().then((res) => {
    console.log(res)
    if(window.timeline?.permalink) {
      let path = window.location.pathname.split("/")[1]
      location.href = `/${path}`
    }
    if(res?.redacted) {
    } else {
    }
  }).then(() => {
  })
}

async function redactPost() {
    let endpoint = `/post/redact`

    let data = {
        room_id: post?.room_id,
        event_id: post?.event_id,
        reason: "dunno",
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

function report() {
  reportPost().then((res) => {
    console.log(res)
    if(res?.redacted) {
    } else {
    }
  }).then(() => {
  })
}

async function reportPost() {
    let endpoint = `/post/report`

    let data = {
        room_id: window.timeline?.room_id,
        event_id: post?.event_id,
        reason: "dunno",
        score: 50,
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

function edit() {
  window.editPost(id)
  menu.hide()
}


let viewing = false;

function viewEvent() {
    menu.hide()
    viewing = true
}

function killViewEvent() {
    viewing = false
}

</script>

<div class="pointer" bind:this={icon}>
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M12.78 6.22a.75.75 0 010 1.06l-4.25 4.25a.75.75 0 01-1.06 0L3.22 7.28a.75.75 0 011.06-1.06L8 9.94l3.72-3.72a.75.75 0 011.06 0z"></path></svg>
</div>

<div class="t-con dis-no" bind:this={content}>

    <div class="post-menu small flex flex-column pa3" style="min-width: 180px">


        {#if showRedact}

        <div class="flex flex-column mb3 hov-bo pointer">
            <span class="" on:click={redact}>Delete Post</span>
        </div>
        {/if}


        {#if authenticated}

        <div class="flex flex-column mb3  hov-bo pointer">
          <span class="" on:click={menu.hide()}>Report Post</span>
        </div>

        <div class="flex flex-column hov-bo pointer">
          <span class="" on:click={viewEvent}>View Source</span>
        </div>

        {/if}

    </div>
</div>

{#if viewing}
<div class="modal-container main ph3" transition:fade="{{duration: 33}}">

  <div class="modal-inner-np start flex flex-column " >

    <div class="pa3 flex flex-column">

      <div class="flex flex-column">
        <span class=" fs-09"><strong>Post Source</strong></span>
      </div>


      <div class="flex mt3 view-box scrl ovfl-y small">
          <pre>
          {JSON.stringify(post, null, 4)}
          </pre>

      </div>



      <div class="flex mt4">

        <div class="flex-one">
        </div>

        <div class="">
          <button class="" on:click={killViewEvent}>Done</button>
        </div>

      </div>




    </div>


  </div>

  <div class="mask absolute" on:click={killViewEvent}></div>

</div>
{/if}


<style>
.post-menu {
    border-radius: 17px;
}
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
    z-index: 50000;

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
    width: 728px;
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


.view-box {

}

</style>


