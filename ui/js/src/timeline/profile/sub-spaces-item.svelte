<script>
export let child;
</script>

<option value={child.path}>{child.path}</option>

{#each child.children as child (child.room_id)}
    <svelte:self child={child} />
{/each}
