# JBrowse Plugins

JBrowse plugins are JavaScript modules that extend the app at runtime. JBrowse fetches them as UMD bundles from URLs listed in `config.json` and calls two lifecycle methods on each plugin class:

```json
"plugins": [
  {
    "name": "HelloWorld",
    "url": "plugins/hello-world.umd.js"
  }
]
```

```js
export default class HelloWorldPlugin {
  name = "HelloWorld"
  version = "1.0.0"

  // Called first — register new types (tracks, widgets, adapters)
  install(pluginManager) {}

  // Called after all plugins are installed — wire up behaviour
  // (menu items, extension points, modifications to existing types)
  configure(pluginManager) {}
}
```

## State models (MobX-State-Tree)

JBrowse uses MobX-State-Tree (MST) for all state. When you register a new type (e.g. a widget), you provide an MST model that describes what data it holds and what actions can change it:

```js
const MyWidgetModel = pluginManager.lib["mobx-state-tree"].types
  .model("MyWidget", {
    id: types.identifier,
    type: types.literal("MyWidget"),
    featureName: types.maybe(types.string),
  })
  .actions((self) => ({
    setFeatureName(name) {
      self.featureName = name
    },
  }))
```

Models have three sections:

| Section                 | Purpose                                           |
| ----------------------- | ------------------------------------------------- |
| `types.model(...)`      | Declares the shape of the data (typed properties) |
| `.actions(self => ...)` | Functions that are allowed to mutate state        |
| `.views(self => ...)`   | Derived/computed values, like getters             |

## React components

React components read from the MST model. Wrapping them in `observer()` makes them re-render automatically whenever observed state changes:

```js
const { observer } = pluginManager.lib["mobx-react"]
const { React } = pluginManager.lib

const MyWidgetComponent = observer(({ model }) => (
  <div>
    {model.featureName ? (
      <p>Selected: {model.featureName}</p>
    ) : (
      <p>Click a feature</p>
    )}
  </div>
))
```

The data flow is:

```
user action
    ↓
MST action mutates state   (e.g. model.setFeatureName("geneA"))
    ↓
MobX detects the change
    ↓
observer() re-renders the component
```

`observer()` is the bridge — it makes the React component subscribe to exactly the state it reads, nothing more.

## Building

Plugin source lives in `plugins/<name>/`. Vite builds each plugin into a UMD bundle written to `jbrowse2/plugins/`, which the static server already serves.

```sh
bun run build:plugins   # one-off build
bun run dev:plugins     # rebuild on save
```

To add a new plugin:

1. Create `plugins/<name>/index.js` with the plugin class as the default export.
2. Add an entry to `vite.plugin.config.ts` pointing at the new entry file.
3. Add the plugin to `config.json` under `"plugins"` with a URL of `plugins/<name>.umd.js`.
4. Run `bun run build:plugins`.
