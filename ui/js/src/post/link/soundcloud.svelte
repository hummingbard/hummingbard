<script>

import {onMount} from 'svelte'


export let title;
export let description;
export let player;
export let image;
export let href;

export let reply;

let active = false;

let iframe;


function activate() {
    active = true
}

</script>

<div class="link-item flex flex-column">

{#if active}
<div class="frame-con">
    <iframe 
        title={title}
        bind:this={iframe}
        src={player}
        width="100%" 
        height="100%" 
        frameborder="0" 
        allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" 
        allowfullscreen>
    </iframe>
</div>
{:else}


  <div class="video-container flex " on:click={activate} >

      <div class="vp-i gr-default bg-img pointer"
        class:vp-i-s={reply}
         style="background-image: url({image});">

        <div class="gr-center">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40" width="40" height="40">
	<defs>
		<image width="40" height="40" id="img1" href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACgAAAAoCAMAAAC7IEhfAAAAAXNSR0IB2cksfwAAAJBQTFRFAAAAAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBYmJidXV1CAgIenp6////5OTkYmJiAQEB7+/vAQEBY2NjdXV1AQEBY2NjAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBsx3KIgAAADB0Uk5TAAw7g77c8/0Yh93/hhmo+qda5Vmcmwe4/tbS8eD/6tT06jrU0fLT/IL5F6bkWJq3xQuJawAAATVJREFUeJyNldl6gyAQhUdFMS7RaozpYtOaWtOs7/92BWsRYTCcO+X/jswwHAFkOa5H/IDSwCee64BJ4SqisRCNkhDF0nUWK8ryVOeeChXjKkoFczYYxlXNtrqtdyZwV28lv9qEcdWT58boN3hWoo4ljGusKEXrlVX8dWktXjy/vOJkPpzH1Oe35n3/gYEZP6NVLIFN84m6Jqw10RxsUNfIAZeqIOZKW/BiHWSuB8XVA4KCmisB3wBy12Ba8iEwgsz1S7h2QBdA5vo9GvWPwIMArT9tXYx1eywbfrQ/QuuhgGQOGsfManB/TnzEc/FsvApny8t1GdOqfARe/y92tRwAN8tIuUsxtRRSdymkmGdlAm9qQpdo7ZcraEpzPZrPeIqHiRL2JxQbttqy30fX951Pju18c7+whEyqDwiKzQAAAABJRU5ErkJggg=="/>
	</defs>
	<style>
		tspan { white-space:pre }
	</style>
	<use id="Background" href="#img1" x="0" y="0" />
</svg>
        </div>
      </div>



  <div class="flex-one ph3 gr-default">
    <div class="flex flex-column gr-start-center">
      <div class="small clmp-2 fs-09">
        <a target="_blank" href={href}>
            <span class="primary">{title}</span>
        </a>
      </div>
      {#if !reply}
      <div class="mt1 smaller o-80 clmp-2">
        <a target="_blank" href={href}>
        {description}
        </a>
      </div>
      {/if}
      <div class="smaller o-80 mt1 clmp-1">
          <a class="underline" target="_blank" href={href}>
            {href}
          </a>
      </div>
    </div>
  </div>



</div>
{/if}

</div>

<style>
.frame-con {
  overflow: hidden;
  padding-top: 56.25%;
  position: relative;
  border-radius: 13px;
  background-color: var(--primary-lightest-gray);
}
 
.frame-con iframe {
   border: 0;
   height: 100%;
   left: 0;
   position: absolute;
   top: 0;
   width: 100%;
}

</style>
