<script>
import { createEventDispatcher } from 'svelte';
const dispatch = createEventDispatcher();

export let store;

$: active = store?.review?.enabled

$: bg = active ? `--primary-dark` : `--primary-grayish`
$: col = active ? `--white` : `--primary-text`

$: tip = active ? `Remove Review` : `Add Review`

function toggle() {
  dispatch('toggleReview')
}

</script>


<div class="pointer o-70 hov-op" class:fill-blue={active} 
     on:click={toggle}
   aria-label={tip}
  data-microtip-position="bottom"
  data-microtip-size="fit"
  role="tooltip">
    {#if !active}
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M8 .25a.75.75 0 01.673.418l1.882 3.815 4.21.612a.75.75 0 01.416 1.279l-3.046 2.97.719 4.192a.75.75 0 01-1.088.791L8 12.347l-3.766 1.98a.75.75 0 01-1.088-.79l.72-4.194L.818 6.374a.75.75 0 01.416-1.28l4.21-.611L7.327.668A.75.75 0 018 .25zm0 2.445L6.615 5.5a.75.75 0 01-.564.41l-3.097.45 2.24 2.184a.75.75 0 01.216.664l-.528 3.084 2.769-1.456a.75.75 0 01.698 0l2.77 1.456-.53-3.084a.75.75 0 01.216-.664l2.24-2.183-3.096-.45a.75.75 0 01-.564-.41L8 2.694v.001z"></path></svg>
    {:else}
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M8 .25a.75.75 0 01.673.418l1.882 3.815 4.21.612a.75.75 0 01.416 1.279l-3.046 2.97.719 4.192a.75.75 0 01-1.088.791L8 12.347l-3.766 1.98a.75.75 0 01-1.088-.79l.72-4.194L.818 6.374a.75.75 0 01.416-1.28l4.21-.611L7.327.668A.75.75 0 018 .25z"></path></svg>
    {/if}
</div>

<style>
.anonymous-button {
  padding: 0.125rem 0.25rem;
  color: black;
  border: 0;
  font-size: small;
  font-weight: bold;
  border-radius: 4px;
}
</style>


