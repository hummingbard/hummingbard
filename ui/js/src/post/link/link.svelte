<script>
import Youtube from './youtube.svelte'
import Vimeo from './vimeo.svelte'
export let link;
export let reply;


</script>


{#if link.is_youtube}
    <Youtube
        id={link.youtube_id}
        title={link.title}
        description={link.description}
        href={link.href}
        reply={reply}
    />

{:else if link.is_vimeo}
    <Vimeo
        id={link.vimeo_id}
        title={link.title}
        description={link.description}
        href={link.href}
        image={link.image}
        reply={reply}
    />

{:else}

<div class="link-item pa3 flex flex-column fs-09">
    <div class="">
      <a href={link.href}>
        <span class="primary">
            {link.title?.length > 0 ? link.title : link.href}
        </span>
      </a>
    </div>
    {#if link.title?.length > 0}
    <div class="small o-90 mt2">
      {link.href}
    </div>
    {/if}
</div>

{/if}
