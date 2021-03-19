<script>
import tippy from 'tippy.js'
import {onMount} from 'svelte'
import { getPostById } from '../../timeline/store.js'
import { u2764, u1f602, u1f92a, u1f62e, u1f914, u1f610, u1f612, u1f614} from '../../emoji/emoji.js'

export let id;
$: post = getPostById(id)

onMount(() => {
      load()
})

let menu;
let content;
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

async function sendReaction() {
    let endpoint = `/post/react`

    let data = {
        room_id: post?.room_id,
        event_id: post?.event_id,
        key: '🥶',
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

function react() {
    sendReaction().then(res => {
        console.log(res)
    }).then(() => {
    })
}

let icon;

</script>

<div class="react pointer ml3" bind:this={icon}>
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M1.5 8a6.5 6.5 0 1113 0 6.5 6.5 0 01-13 0zM8 0a8 8 0 100 16A8 8 0 008 0zM5 8a1 1 0 100-2 1 1 0 000 2zm7-1a1 1 0 11-2 0 1 1 0 012 0zM5.32 9.636a.75.75 0 011.038.175l.007.009c.103.118.22.222.35.31.264.178.683.37 1.285.37.602 0 1.02-.192 1.285-.371.13-.088.247-.192.35-.31l.007-.008a.75.75 0 111.222.87l-.614-.431c.614.43.614.431.613.431v.001l-.001.002-.002.003-.005.007-.014.019a1.984 1.984 0 01-.184.213c-.16.166-.338.316-.53.445-.63.418-1.37.638-2.127.629-.946 0-1.652-.308-2.126-.63a3.32 3.32 0 01-.715-.657l-.014-.02-.005-.006-.002-.003v-.002h-.001l.613-.432-.614.43a.75.75 0 01.183-1.044h.001z"></path></svg>
</div>

<div class="t-con dis-no" bind:this={content}>

    <div class="post-menu small flex  pa2">

        <div class="emo mr2 gr-default">
            <div class="emoji gr-center">
                {@html u2764}
            </div>
        </div>


        <div class="emo mr2 gr-default">
            <div class="emoji gr-center">
                {@html u1f602}
            </div>
        </div>

        <div class="emo mr2 gr-default">
            <div class="emoji gr-center">
                {@html u1f92a}
            </div>
        </div>

        <div class="emo mr2 gr-default">
            <div class="emoji gr-center">
                {@html u1f62e}
            </div>
        </div>

        <div class="emo mr2 gr-default">
            <div class="emoji gr-center">
                {@html u1f914}
            </div>
        </div>

        <div class="emo mr2 gr-default">
            <div class="emoji gr-center">
                {@html u1f610}
            </div>
        </div>

        <div class="emo mr2 gr-default">
            <div class="emoji gr-center">
                {@html u1f612}
            </div>
        </div>

        <div class="emo mr2 gr-default">
            <div class="emoji gr-center">
                {@html u1f614}
            </div>
        </div>


    </div>
</div>

<style>
.emo {
    width: 28px;
    height: 28px;
}
.emoji {
    width: 22px;
    height: 22px;
    cursor: pointer;
    transition: 0.06s;
    transition-timing-function: ease;
}
.emoji:hover {
    width: 28px;
    height: 28px;
}
</style>
