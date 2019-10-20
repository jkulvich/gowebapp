class Api {
    constructor(addr) {
        this.addr = addr;
    }

    async store(data) {
        let resp = await fetch(`${this.addr}/data`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
        if (!resp.ok) throw 'store not OK error: ' + JSON.stringify(resp);
    }

    async load() {
        let resp = await fetch(`${this.addr}/data`);
        if (!resp.ok) throw 'load not OK error: ' + JSON.stringify(resp);
        return await resp.json()
    }
}

class App {
    constructor(api) {
        this.api = api;

        let toolbar = document.querySelector('#toolbar');
        let btnLoad = toolbar.querySelector('#btn-load');
        let btnStore = toolbar.querySelector('#btn-store');
        this.toolbar = {btnLoad, btnStore};

        this.toolbar.btnLoad.onclick = this.loadDataHandler.bind(this);
        this.toolbar.btnStore.onclick = this.storeDataHandler.bind(this);

        let itemNew = document.querySelector('#item-new');
        let btnNew = itemNew.querySelector('#btn-new');
        let inputKey = itemNew.querySelector('#input-key');
        let inputValue = itemNew.querySelector('#input-value');
        this.itemNew = {itemNew, btnNew, inputKey, inputValue};

        this.itemsContainer = document.querySelector('#items-container');

        this.itemNew.btnNew.onclick = this.newItemHandler.bind(this);

        this.itemTemplate = document.querySelector('#item-template');
    }

    createItem(key, value) {
        let item = document.importNode(this.itemTemplate.content, true);
        let inputKey = item.querySelector('#input-key');
        let inputValue = item.querySelector('#input-value');
        let btnRemove = item.querySelector('#btn-remove');
        inputKey.value = key;
        inputValue.value = value;
        let child = item.children[0];
        window.w = child;
        btnRemove.onclick = () => this.removeItemHandler(child);
        return item;
    }

    exportItems() {
        let nodes = this.itemsContainer.querySelectorAll('.item');
        let items = [];
        nodes.forEach(n => {
            items.push({
                key: n.querySelector('#input-key').value,
                value: n.querySelector('#input-value').value,
            });
        });
        return items;
    }

    importItems(data) {
        data.forEach(d => {
            let item = this.createItem(d.key, d.value);
            this.itemsContainer.insertBefore(item, this.itemNew.itemNew);
        });
    }

    clearItems() {
        let nodes = this.itemsContainer.querySelectorAll('.item');
        nodes.forEach(n => this.itemsContainer.removeChild(n));
    }

    async storeDataHandler() {
        let data = this.exportItems();
        this.api.store(data);
    }

    async loadDataHandler() {
        let data = await this.api.load();
        this.clearItems();
        this.importItems(data);
    }

    async newItemHandler() {
        let key = this.itemNew.inputKey.value;
        let value = this.itemNew.inputValue.value;
        let item = this.createItem(key, value);
        this.itemNew.inputKey.value = '';
        this.itemNew.inputValue.value = '';
        this.itemsContainer.insertBefore(item, this.itemNew.itemNew);
    }

    async removeItemHandler(item) {
        this.itemsContainer.removeChild(item);
    }
}

window.onload = async () => {
    let api = new Api('/api');
    let app = new App(api);
    app.loadDataHandler();
};