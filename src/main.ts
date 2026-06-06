import { mount } from 'svelte';
import App from './App.svelte';
import './app.css';

const appEl = document.getElementById('app')!;
appEl.innerHTML = ''; // Remove static HTML startup loader before Svelte mounts
const app = mount(App, { target: appEl });

export default app;
