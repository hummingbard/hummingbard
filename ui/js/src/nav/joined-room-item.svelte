<script>
import { onMount, createEventDispatcher } from 'svelte'
import tippy from 'tippy.js'
const dispatch = createEventDispatcher()
export let room;

function alias(alias) {
    let x = alias.substring(1)
    x = x.replace(`:${window.location.hostname}`, "")
    x = x.replaceAll("_", "/")
    return x
}


let options;
let content;
let menu;

onMount(() => {
      load()
})

let highlight = false;

function load() {
    menu = tippy(options, {
        content: content,
        theme: 'options',
        placement: "bottom",
        duration: "40",
        animation: "shift-away" ,
        trigger: "click",
        arrow: false,
        interactive: true,
        onHide(menu) {
            highlight = false
        },
        onShow(menu) {
            highlight = true
        },
    });

    content.classList.remove('dis-no')
}

async function leave() {
    let endpoint = `/room/leave`

    let data = {
        id: room.room_id,
        alias: room.room_alias,
    }

    let options = {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Authorization': identity.access_token,
            'Content-Type': 'application/json'
        }
    }
    console.log(options)

    let resp = await fetch(endpoint, options)
    const ret = await resp.json()
    return Promise.resolve(ret)
}


function leaveRoom() {
    menu.hide()
  leave().then((res) => {
    console.log(res)
      if(res?.left) {
          window.navUpdateJoinedRoom('left', room)
      }
  }).then(() => {
  })
}

</script>

<a href="/{alias(room.room_alias)}">
    <div class="pv1 flex room-item ph1 jr-i" class:high={highlight}>
    {#if room.avatar?.length > 0}
        <div class="thumbnail-s mr3 gr-center">
          <img src="{room.avatar}" />
        </div>
    {:else}
        <div class="mr3 gr-center">
            <svg class="gr-center" height="22" width="22">
               <circle cx="11" cy="11" r="11" stroke-width="0" fill="black" />
            </svg>
        </div>
    {/if}
    <div class="gr-center flex-one ov-x-h flex" title={room.room_alias}>
        <div class="primary small ov-x-h t-ov-h mr3 flex-one gr-center">{alias(room.room_alias)}</div>
        <div class="gr-center pa1 jr-m" 
            class:jr-h={highlight}
            title="Space Options"
            bind:this={options}
            on:click|preventDefault="">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path d="M4.427 7.427l3.396 3.396a.25.25 0 00.354 0l3.396-3.396A.25.25 0 0011.396 7H4.604a.25.25 0 00-.177.427z"></path></svg>
        </div>
    </div>
</div>
</a>

<div class="t-con dis-no" bind:this={content}>

    <div class="post-menu small flex flex-column" style="min-width: 180px">



        <div class="flex flex-column pointer op-i ph3 pv2 st"
             on:click|preventDefault={leaveRoom}>
            <span class="">Leave Room</span>
        </div>


    </div>
</div>
<style>
.room-item {
    border-radius: 500px;
}
.room-item:hover {
    background: var(--primary-grayish);
}

.high {
    background: var(--primary-grayish);
    opacity: 1;
}

.jr-i:hover .jr-m {
    opacity: 1;
}
.jr-m {
    opacity: 0;
}
.jr-h {
    opacity: 1;
}

.op-i:hover {
    background: var(--primary-light-gray);
}
.st {
    border-radius: 13px;
}
</style>
