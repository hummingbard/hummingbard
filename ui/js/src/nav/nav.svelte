<script>
import {fade} from 'svelte/transition'
import {onMount} from 'svelte'
import ExploreSpaces from './explore-spaces.svelte'
import JoinedRooms from './joined-rooms.svelte'

let active = false;



onMount(() => {
    /*
    let n = localStorage.getItem("nav");
    if(n && n == 'true') {
        active = true
        let el = document.querySelector('.nav-de')
        if(el) {
            el.classList.add('nav-de-ac')
        }
    }
    */
    window.exploreSpaces = (ok) => {
        active = true
        let el = document.querySelector('.nav-de')
        if(el) {
            el.classList.add('nav-de-ac')
        }
        exploring = true
    }
})

function toggle() {
    active = !active
    if(active) {
        let el = document.querySelector('.nav-de')
        if(el) {
            el.classList.add('nav-de-ac')
        }
    } else {
        let el = document.querySelector('.nav-de')
        if(el) {
            el.classList.remove('nav-de-ac')
        }
    }
    localStorage.setItem("nav", active);
}


let exploring = false;

function explore() {
    exploring = !exploring
}


</script>

<div class="n flex flex-column">

    <div class="hg gr-default pointer" 
         style="min-height: 56px;" 
              aria-label="Open Sidebar"
              data-microtip-position="right"
              data-microtip-size="fit"
              role="tooltip"
        on:click={toggle}>
        <div class="gr-center pv3" >
            <svg class="hgg" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M1 2.75A.75.75 0 011.75 2h12.5a.75.75 0 110 1.5H1.75A.75.75 0 011 2.75zm0 5A.75.75 0 011.75 7h12.5a.75.75 0 110 1.5H1.75A.75.75 0 011 7.75zM1.75 12a.75.75 0 100 1.5h12.5a.75.75 0 100-1.5H1.75z"></path></svg>
        </div>
    </div>

</div>

{#if active}
    <div class="nav flex flex-column"
    transition:fade="{{duration: 53}}">

        <div class="flex fl-o pa3">

            <div class="gr-center flex-one mr3">
                <button class="light" style="width: 100%;" on:click={explore}>{exploring ? 'Done' : 'Explore Spaces'}</button>
            </div>


            <div class="gr-default pointer" on:click={toggle}>
                <div class="gr-center ">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M15.28 5.22a.75.75 0 00-1.06 0l-6.25 6.25a.75.75 0 000 1.06l6.25 6.25a.75.75 0 101.06-1.06L9.56 12l5.72-5.72a.75.75 0 000-1.06z"></path></svg>
                </div>
            </div>

        </div>



        <div class="n-contain h-100 flex flex-column scr ">
            {#if exploring}
                <ExploreSpaces 
                on:kill={explore}/>
            {/if}


            <JoinedRooms />

        </div>


    </div>
{/if}

<style>
.n {
    width: 100%;
    position: sticky;
    top: 0;
    height: 100vh;
}
.n-itm {
}
.nav {
  position :fixed;
  width: 300px;
  background: var(--primary-lightestest-gray);
  top: 0;
  bottom: 0;
  left: 0;
}

.n-contain {
    overflow-x: auto;
}


.scr  {
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: var(--primary-gray) var(--primary-light-gray);
}

.scr::-webkit-scrollbar {
  width: 6px;
}
.scr::-webkit-scrollbar-track {
  background: var(--primary-light-gray);
}
.scr::-webkit-scrollbar-thumb {
  background-color: var(--primary-gray);
}

.hg:hover .hgg{
    stroke-width: 3;
    fill: red;
}

@media screen and (max-width: 768px) {
    .n {
        width: 2rem;
    }
}


</style>
