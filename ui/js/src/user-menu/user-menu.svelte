<script>
import tippy from 'tippy.js';
import {onMount} from 'svelte'

$: avatarExists = identity?.avatar_url.length > 0
$: avatarURL = identity?.avatar_url


$: ident = identity

let content;
let icon;
let menu;

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

function username() {
    let server = identity?.user_id?.split(":")[1]
    if(server == window.location.hostname) {
        return identity?.user_id?.split(":")[0]
    }
    return identity?.user_id
}
</script>

<div bind:this={icon} class="thumbnail pointer">
{#if avatarExists}
    <img alt={ident?.user_id} src={avatarURL} />
{:else}
   <svg class="gr-center" height="30" width="30">
     <circle cx="15" cy="15" r="15" stroke-width="0" fill="black" />
   </svg>
{/if}
</div>


<div class="t-con dis-no" bind:this={content}>

    <div class="header-menu flex flex-column fs-09" style="min-width: 180px">

        <div class="flex flex-column ph3 pv2  hov-bo">
            <a href="/{username()}">
                <span class="">Profile ({username()})</span>
            </a>
        </div>

        <div class="flex flex-column ph3 pv2  hov-bo">
            <a href="/invite">
                <span class="">Invite A Friend</span>
            </a>
        </div>

        <div class="flex flex-column ph3 pv2  hov-bo">
            <a href="/settings">
                <span class="">Account Settings</span>
            </a>
        </div>

        <div class="flex flex-column ph3 pv2  hov-bo">
            <a href="/logout">
                <span class="">Logout</span>
            </a>
        </div>

    </div>
</div>


