<script>
import {onMount} from 'svelte'
import { getPostById, getPostReactions } from '../../timeline/store.js'
import ReactionItem from './reaction-item.svelte'
import React from '../react/react.svelte'

export let id;

$: post = getPostById(id)

$: reactions = getPostReactions(id)

let items = [];

onMount(() => {
    items = categorize(reactions)
})

function categorize(reactions) {
    let items = []
    reactions.forEach(reaction => {
        let ind = items.findIndex(x => x.key == reaction.content.m_relates_to.key)
        if(ind == -1) {
            items.push({
                key: reaction.content.m_relates_to.key,
                senders: [{event_id: reaction.event_id, sender: reaction.sender}]
            })
        } else {
            items[ind].senders.push({event_id: reaction.event_id, sender: reaction.sender})
        }
    })
    items.sort((a, b) => (a.senders.length < b.senders.length) ? 1 : -1)
    return items
}


function reacted(e) {
    console.log(e.detail)
    /*
    items.forEach(reaction => {
        let match = reaction.key == e.detail
        let ind = reaction.senders.findIndex(x => x.sender == identity.user_id)
        if(match && (ind != -1)) {
            reaction.senders.splice(ind, 1)
            items = items
        }
        if(reaction.senders.length == 0) {
            let ind = items.findIndex(x => x.key = e.detail)
            if(ind != -1) {
                items.splice(ind, 1)
                items = items
            }
        }
    })
    */
    let ind = items.findIndex(x => x.key == e.detail)
    if(ind == -1) {
        sendReaction(e.detail).then(res => {
            console.log(res)
            if(res?.event) {
                let item = {
                    key: e.detail,
                    senders: [{event_id: res?.event?.event_id, sender: identity.user_id}]
                }
                items = [...items, item]
            }
        }).then(() => {
        })
    } else {
        let item = items[ind].senders.filter(x => x.sender == identity.user_id)
        let reacted = item.length == 1
        if(!reacted) {
            sendReaction(e.detail).then(res => {
                console.log(res)
                if(res?.event) {
                    let item = {
                        key: e.detail,
                        senders: [{event_id: res?.event?.event_id, sender: identity.user_id}]
                    }
                    items[ind].senders.push({event_id: res?.event?.event_id, sender: identity.user_id})
                    items = [...items, item]
                }
            }).then(() => {
            })
        } else {
            redactReaction(item[0].event_id).then(res => {
                console.log(res)
                if(res?.redacted) {
                    let i = items.findIndex(x => x.key == e.detail)
                    if(i != -1) {
                        let se = items[i].senders.findIndex(x => x.sender = identity.user_id)
                        if(se != -1) {
                            items[i].senders.splice(se, 1)
                            items = items
                            clearEmpty()
                        }
                    }
                }
            }).then(() => {
            })
        }
    }
}

function clearEmpty() {
    items.forEach(item => {
        if(item.senders.length == 0) {
            let ind = items.findIndex(x => x.key == item.key)
            if(ind != -1) {
                items.splice(ind, 1)
                items = items
            }
        }
    })
}


async function sendReaction(key) {
    let endpoint = `/post/react`
    let data = {
        room_id: post?.room_id,
        event_id: post?.event_id,
        key: key,
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

async function redactReaction(id) {
    let endpoint = `/post/redact`
    let data = {
        room_id: post?.room_id,
        event_id: id,
        reason: "remove reaction",
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

</script>

<div class="flex">

    <div class="flex mr3 gr-center">
        <React id={id} on:reacted={reacted}/>
    </div>

    <div class="flex flex-wrap">
    {#each items as item (item.key)}
        <div class="mr2">
            <ReactionItem 
             item={item} 
             id={id} 
             on:reacted={reacted}/>
        </div>
    {/each}
    </div>

</div>

