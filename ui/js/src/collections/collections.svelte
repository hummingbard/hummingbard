<script>
import {onMount} from 'svelte'
import { initializeCollections, storeCollections } from './store.js'
import AddCollection from './add-collection.svelte'

onMount(() => {
    initializeCollections()
})

$: exists = $storeCollections?.length > 0

$: owner = window.timeline?.owner && authenticated
$: user = window.timeline?.room_path

let adding = false;

function add() {
    adding = true
}

function killAdd() {
    adding = false
}

</script>

<div class="flex flex-column">


{#if adding && owner}
    <AddCollection on:kill={killAdd}/>
{/if}

{#if exists}

    {#if !adding && owner}
    <div class="brd-btm pa3">
        <button class="" on:click={add}>Add a collection</button>
    </div>
    {/if}

    <div class="flex flex-column flex-one">
        {#each $storeCollections as collection (collection.id)}
            <a href="/{collection.path}">
            <div class="mb3 pa3 flex flex-column brd-btm">
                <div class="">
                    <strong>{collection.name}</strong>
                </div>
                <div class="mt3 fs-09">
                    {collection.description}
                </div>
            </div>
            </a>
        {/each}
    </div>

{:else}
    {#if !adding}
        <div class="gr-default h-100">
            <div class="gr-center flex flex-column">
                {#if owner}
                    <div class="tc lh-copy">
                        <span class="fs-09">You don't have any collections yet.</span>
                    </div>
                    <div class="mt3 tc">
                        <button class="" on:click={add}>Add a collection</button>
                    </div>
                {:else}
                    <div class="tc">
                        <span class="fs-09"><a href="/{user}"><span class="primary">{user}<span></a> doesn't have any collections yet.</span>
                    </div>
                {/if}
            </div>
        </div>
    {/if}
{/if}


</div>
