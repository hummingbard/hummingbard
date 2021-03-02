<script>
import Join from '../join/join.svelte'
export let space;

function alias(alias) {
    let x = alias?.substring(1)
    x = x.replace(`:${window.location.hostname}`, "")
    x = x.replaceAll("_", "/")
    return x
}

$: spaceAlias = alias(space?.room_alias)
</script>

<div class="mb3 item pa2 flex flex-column mr3">
  <div class="flex">
    {#if space.avatar?.length > 0}
        <div class="thumbnail mr3">
          <img src="{space.avatar}" />
        </div>
    {:else}
        <div class="mr3">
            <svg class="gr-center" height="33" width="33">
               <circle cx="16.5" cy="16.5" r="16.5" stroke-width="0" fill="black" />
            </svg>
        </div>
    {/if}
    <div class="flex flex-column flex-one">
      <div class="mt2">
        <strong>{space.name}</strong>
      </div>
      {#if space.topic.length > 0}
        <div class="mt2 lh-copy no-m small">
          {@html space.topic}
        </div>
      {/if}
      <div class="mt2">
        <a href="/{space.room_path}"><span class="small primary">{space.room_path}</span></a>
      </div>
    </div>
    <div class="">
      <Join
        type="space"
        alias={space.room_alias}
        id={space.room_id}
        inline={true}
        />
    </div>
  </div>
</div>

<style>
.item {
  border: 1px solid var(--primary-grayish);
  border-radius: 13px;
}
.topic p{
  margin: 0;
}
</style>
