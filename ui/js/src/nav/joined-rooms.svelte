<script>
import {tick} from 'svelte'
import {fade} from 'svelte/transition'
import JoinedRoomItem from './joined-room-item.svelte'

let joined_rooms = [];
if(identity?.joined_rooms) {
    identity?.joined_rooms.forEach(room => {
        if(!room.room_alias.includes('@')) {
            let ind = joined_rooms.findIndex(x => x.room_alias == room.room_alias)
            if(ind = -1) {
                joined_rooms.push(room)
            }
        }
    })
    joined_rooms?.sort((a, b) => (a.room_alias > b.room_alias) ? 1 : -1)
}

let limit = 13

function alias(alias) {
    let x = alias.substring(1)
    x = x.replace(`:${window.location.hostname}`, "")
    x = x.replaceAll("_", "/")
    return x
}
function filtered() {
    return joined_rooms.filter(x => alias(x.room_alias).includes(query))
}

$: shortened = (searching && query.length > 0) ? filtered() : (searching &&
    query.length == 0) ? joined_rooms : showingMore ? joined_rooms : joined_rooms.slice(0, limit);

$: count = joined_rooms.length

let active = true

function toggle() {
    active = !active
}

$: showMore = joined_rooms.length > limit

let showingMore = false;

function showAll() {
    showingMore = !showingMore
}

let searching = false;
let searchInput;
let query = '';

function updateQuery(e) {
    if(e.key == 'Escape') {
        query = ''
        searching = false
    }
}

async function toggleSearch() {
    searching = !searching
    if(!searching) {
        searchInput.value = null
    } else {
        await tick();
        searchInput.focus()
    }
}


</script>

    <div class="flex ph3 mt2">
        <div class="mr2 pointer gr-center" on:click={toggle}>
            {#if active}
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M12.78 6.22a.75.75 0 010 1.06l-4.25 4.25a.75.75 0 01-1.06 0L3.22 7.28a.75.75 0 011.06-1.06L8 9.94l3.72-3.72a.75.75 0 011.06 0z"></path></svg>
            {:else}
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M3.22 9.78a.75.75 0 010-1.06l4.25-4.25a.75.75 0 011.06 0l4.25 4.25a.75.75 0 01-1.06 1.06L8 6.06 4.28 9.78a.75.75 0 01-1.06 0z"></path></svg>
            {/if}
        </div>
        <div class="gr-center pointer" on:click={toggle}>
            <span class="bold small no-select">Joined Spaces</span>
        </div>
        {#if joined_rooms.length > 0}
        <div class="gr-default count ml2">
            <div class="gr-center">
                <span class="smallest bold "> {joined_rooms.length > 0 ? count : ''}</span>
            </div>
        </div>
        {/if}
        <div class="flex-one"></div>
        {#if active}
            <div class="gr-center ml2 pointer" on:click={toggleSearch}>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M11.5 7a4.499 4.499 0 11-8.998 0A4.499 4.499 0 0111.5 7zm-.82 4.74a6 6 0 111.06-1.06l3.04 3.04a.75.75 0 11-1.06 1.06l-3.04-3.04z"></path></svg>
        </div>
        {/if}
    </div>
    {#if active}

    {#if searching}
        <div class="ph3 mt2 small">
            <input bind:this={searchInput} bind:value={query}
                   on:keydown={updateQuery}
            placeholder="Filter Joined Spaces"/>
        </div>
    {/if}

    <div class="shm mt3 ph2 lh-copy" transition:fade="{{duration: 53}}">
        {#if showMore || searching}
            {#each shortened as room (room.room_id)}
                <JoinedRoomItem room={room}/>
            {/each}
        {:else}
            {#each joined_rooms as room (room.room_id)}
                <JoinedRoomItem room={room}/>
            {/each}
        {/if}

        {#if searching && shortened.length == 0}
            <div class="gr-default pa2 small ">
                    No matches.
            </div>
        {/if}

        {#if joined_rooms.length == 0}
            <div class="fs-09 ph2">
                You have't joined any spaces yet.
            </div>
        {/if}
    </div>
    {/if}
    {#if showMore && active && !searching}
        <div class="shm flex flex-column pv2 ph3 pointer" on:click={showAll}>
            <div class="flex">
                <div class="mr2 pointer gr-center">
                    {#if showingMore}
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M3.22 9.78a.75.75 0 010-1.06l4.25-4.25a.75.75 0 011.06 0l4.25 4.25a.75.75 0 01-1.06 1.06L8 6.06 4.28 9.78a.75.75 0 01-1.06 0z"></path></svg>
                    {:else}
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M12.78 6.22a.75.75 0 010 1.06l-4.25 4.25a.75.75 0 01-1.06 0L3.22 7.28a.75.75 0 011.06-1.06L8 9.94l3.72-3.72a.75.75 0 011.06 0z"></path></svg>
                    {/if}
                </div>
                <div class="gr-center pointer">
                    <span class="bold smaller no-select">{showingMore ? 'Show Less' : 'Show More'}</span>
                </div>
            </div>
            <div class="gr-default ">
                <div class="show-hint gr-center">
                </div>
            </div>
        </div>
    {/if}

<style>
.count {
    background: var(--primary-dark);
    color: var(--white);
    border-radius: 500px;
    padding: 0.125rem 0.5rem;
}

.shm:hover .show-hint {
    opacity: 1;
}
.show-hint {
    transition: 0.1s;
    opacity: 0;
    width: 3rem;
    height: 0.25rem;
    background: var(--primary-gray);
    border-radius: 500px;
}
</style>
