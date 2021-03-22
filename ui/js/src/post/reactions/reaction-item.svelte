<script>
import tippy from 'tippy.js'
import {onMount, createEventDispatcher} from 'svelte'
const dispatch = createEventDispatcher();
import {getPostById, addPost} from '../../timeline/store.js'
import { u2764, u1f602, u1f92a, u1f62e, u1f914, u1f610, u1f612, u1f614} from '../../emoji/emoji.js'
export let item;
export let id;
$: post = getPostById(id)

onMount(() => {
    if(isMember) {
      load()
    }
})

let menu;
let content;
let icon;

function load() {
    menu = tippy(icon, {
        content: content,
        theme: 'menu',
        placement: "right",
        duration: "40",
        animation: "shift-away" ,
        arrow: true,
        onHide(menu) {
        },
        onShow(menu) {
        },
    });

    content.classList.remove('dis-no')
}

function react() {
    if(isMember) {
        dispatch('reacted', item.key)
    }
}

$: reacted = item.senders.filter(x => x.sender == identity.user_id).length > 0

function formatUsername(username) {
    if(username.includes(window.location.hostname)) {
        let x = username.split(":")
        return x[0]
    }
    return username
}

$: isMember = window.timeline?.member

</script>

<div class="co-item flex no-select" bind:this={icon} class:reacted={reacted} on:click={react}>

    {#if item.key == "❤"}
        <div class="emo gr-default">
            <div class="emoji gr-center">
                <div class="no-click">{@html u2764}</div>
            </div>
        </div>
    {:else if item.key == "😂"}
        <div class="emo gr-default">
            <div class="emoji gr-center">
                <div class="no-click">{@html u1f602}</div>
            </div>
        </div>
    {:else if item.key == "🤪"}
        <div class="emo gr-default">
            <div class="emoji gr-center">
                <div class="no-click">{@html u1f92a}</div>
            </div>
        </div>
    {:else if item.key == "😮"}
        <div class="emo gr-default">
            <div class="emoji gr-center">
                <div class="no-click">{@html u1f62e}</div>
            </div>
        </div>
    {:else if item.key == "🤔"}
        <div class="emo gr-default">
            <div class="emoji gr-center">
                <div class="no-click">{@html u1f914}</div>
            </div>
        </div>
    {:else if item.key == "😐"}
        <div class="emo gr-default">
            <div class="emoji gr-center">
                <div class="no-click">{@html u1f610}</div>
            </div>
        </div>
    {:else if item.key == "😒"}
        <div class="emo gr-default">
            <div class="emoji gr-center">
                <div class="no-click">{@html u1f612}</div>
            </div>
        </div>
    {:else if item.key == "😔"}
        <div class="emo gr-default">
            <div class="emoji gr-center">
                <div class="no-click">{@html u1f614}</div>
            </div>
        </div>
    {/if}

    <div class="mh1 small gr-center o-90">
        {item.senders.length}
    </div>
</div>

<div class="r-c dis-no" bind:this={content}>

    <div class="smaller ph2 pv1 flex flex-column">
        {#each item.senders as sender (sender.sender)}
            <div class="pv1">
                <strong>{formatUsername(sender.sender)} </strong>
            </div>
        {/each}

    </div>

</div>

<style>
.r-c {
    background: var(--primary-dark);
    color: var(--white);
    border-radius: 7px;
}

.co-item {
    border: 1px solid var(--primary-lightest-gray);
    border-radius: 500px;
    padding: 0 0.125rem;
    cursor: pointer;
    transition: 0.1s;
}
.co-item:hover {
    border: 1px solid var(--primary-dark-gray);
}
.reacted {
    background: var(--primary-lightest-gray);
}
.emo {
    width: 22px;
    height: 22px;
}
.emoji {
    width: 16px;
    height: 16px;
}
</style>
