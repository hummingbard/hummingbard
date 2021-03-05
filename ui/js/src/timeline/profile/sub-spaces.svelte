<script>
import SubSpacesItem from './sub-spaces-item.svelte'
let selected;
function switchToSpace() {
    if(selected == '') {
        return
    } else{
        window.location.href = `${window.location.pathname}/${selected}`
    }
}

$: params = new URLSearchParams(window.location.search)
$: paramsSort = params.get('sort')
$: recent = paramsSort == 'recent'
$: replies = paramsSort == 'replies'

$: children = window.timeline?.children

</script>

<select class="small" bind:value={selected} on:change={switchToSpace}>
    <option>Sub-Spaces</option>
    <option disabled>----------</option>
    {#each children as child (child.room_id)}
        <SubSpacesItem child={child} />
    {/each}
</select>
