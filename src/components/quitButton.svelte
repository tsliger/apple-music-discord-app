<script lang="ts">
    import { exit } from "@tauri-apps/plugin-process";
    import { invoke } from "@tauri-apps/api/core";

    const quitApp = async () => {
        await invoke("call_kill_api", {});
        await exit(1);
    };
</script>

<button class="quit-button" onclick={quitApp} aria-label="Quit">Quit</button>

<style>
    :global(button) {
        background: none;
        color: inherit;
        border: none;
        padding: 0;
        font: inherit;
        cursor: pointer;
        outline: inherit;
    }

    .quit-button {
        background-color: red;
        padding-top: 10px;
        padding-bottom: 10px;
        width: 64px;
        border-radius: 24px;
        transition: background-color 200ms;
        color: whitesmoke;
        font-weight: 600;
        font-family: Arial, Helvetica, sans-serif;
    }

    .quit-button:hover {
        background-color: rgb(196, 0, 0);
    }
</style>
