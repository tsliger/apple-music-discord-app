<script>
    import ControlPanel from "../components/controlPanel.svelte";
    import { onMount } from "svelte";

    // Prevent right click context menu from appearing
    onMount(() => {
        document.addEventListener("contextmenu", function (e) {
            e.preventDefault();
        });
    });
</script>

<main class="container">
    <ControlPanel />
</main>

<style>
    :global(:root, html, body) {
        margin: 0;
        padding: 0;
        font-family:
            system-ui,
            -apple-system,
            BlinkMacSystemFont,
            "Segoe UI",
            Roboto,
            Oxygen,
            Ubuntu,
            Cantarell,
            "Open Sans",
            "Helvetica Neue",
            sans-serif;
        -webkit-user-select: none;
        -ms-user-select: none;
        user-select: none;
        cursor: default;
    }

    .container {
        --container-padding: 20px;
        padding: var(--container-padding);
        height: calc(100vh - (var(--container-padding) * 2));
        width: calc(100% - (var(--container-padding) * 2));
        display: flex;
        flex-direction: column;
        color: whitesmoke;
    }
</style>
