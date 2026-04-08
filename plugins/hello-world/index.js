export default class HelloWorldPlugin {
  name = "HelloWorld"
  version = "1.0.0"

  install(pluginManager) {
    console.log("[HelloWorld] install()")
  }

  configure(pluginManager) {
    console.log("[HelloWorld] configure()")
  }
}
