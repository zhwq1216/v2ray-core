<template>
    <div id="app">
        <el-tabs type="border-card" v-model="tagValue" @tab-click="tabChanged" ref="tabs">
            <el-tab-pane label="参数配置" name="config">
                <Home/>
            </el-tab-pane>
            <el-tab-pane label="PAC配置" name="pac">
                <PacConfig/>
            </el-tab-pane>
            <el-tab-pane label="流量监控" name="trafficMonitor">
                <TrafficMonitor data-tab-content-name="trafficMonitor" ref="trafficMonitor"/>
            </el-tab-pane>
            <el-tab-pane label="访问日志" name="accessLog">
                <LogViewer log-type="access" data-tab-content-name="accessLog"  ref="accessLog"/>
            </el-tab-pane>
            <el-tab-pane label="普通日志" name="errorLog">
                <LogViewer log-type="error" data-tab-content-name="errorLog" ref="errorLog"/>
            </el-tab-pane>
        </el-tabs>
    </div>
</template>
<script>
    import Home from "@/views/Home";
    import TrafficMonitor from "@/views/TrafficMonitor";
    import LogViewer from "@/components/LogViewer";
    import PacConfig from "@/views/PacConfig";

    export default {
        name: 'App',
        components: {
            Home, TrafficMonitor, LogViewer,PacConfig
        },
        data() {
            return {
                tagValue: "config"
            }
        },
        methods: {
            tabChanged() {
                this.$refs.tabs.panes.forEach((v)=>{
                    let tabContentVue = v.$children[0]
                    if(typeof(tabContentVue.pauseRefresh)==="function") {
                        tabContentVue.pauseRefresh(!v.active)
                    }
                })
            }
        }
    }
</script>
<style>
    #app {
        font-family: Avenir, Helvetica, Arial, sans-serif;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
        color: #2c3e50;
    }

    #nav {
        padding: 30px;
    }

    #nav a {
        font-weight: bold;
        color: #2c3e50;
    }

    #nav a.router-link-exact-active {
        color: #42b983;
    }
</style>
