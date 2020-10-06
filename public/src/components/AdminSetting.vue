<template>
    <setting-card title="Admin管理台设置" :show-enable="true" @update:enableSetting="enableSettingChanged"
                  :enableSetting.sync="enableSetting">
        <el-form ref="form" label-width="80px">
            <el-form-item label="监听端口">
                <template v-slot:label>
                    <el-tooltip>
                        <div slot="content">监听端口可以只配置端口或者指定为监听ip+端口形式</div>
                        <label>监听端口</label>
                    </el-tooltip>
                </template>
                <el-input v-model="sForm.addr" v-setting></el-input>
            </el-form-item>
            <el-form-item label="contextPath">
                <template v-slot:label>
                    <el-tooltip >
                        <div slot="content" >admin web前缀路径</div>
                        <label>contextPath</label>
                    </el-tooltip>
                </template>
                <el-input v-model="sForm.contextPath" v-setting></el-input>
            </el-form-item>
            <el-form-item label="publicPath">
                <el-input v-model="sForm.publicPath" placeholder="应用静态资源文件路径" v-setting></el-input>
            </el-form-item>
        </el-form>
    </setting-card>

</template>

<script>

    export default {
        name: "AdminSetting",
        model:{
            prop:"setting",
            event:"change"
        },
        data() {
            return {
                "enableSetting": true,
                changedByForm: false,
                sForm: {
                    "addr":":8089",
                    "contextPath": "/v2ray",
                    "publicPath": "../public/v2ray/dist"
                }
            }
        },
        created() {
            this.fillDefaultValue(this.setting);
        },
        mounted() {

        },
        methods: {
            enableSettingChanged() {
                this.$nextTick(() => {
                    this.formChanged();
                });
            },
            fillDefaultValue(setting){
                setting = setting || {};
                Object.assign(this.sForm, setting);
                this.$nextTick().then(()=>{
                    this.formChanged();
                });
            },
            formChanged() {
                let setting = this.getSettings();
                if(this.setting !== setting) {
                    this.changedByForm = true;
                    this.$emit("change", setting);
                }

            },
            getSettings() {
                return Object.assign({},this.sForm);
            },
        },
        watch: {
            setting: {
                handler: function (val) {
                    if(this.changedByForm){
                        this.changedByForm = false;
                        return;
                    }
                    this.fillDefaultValue(val);
                },
                deep: false
            }
        },
        props: {
            setting: {
                type: Object,
            }
        }
    }
</script>

<style scoped>
    .el-select {
        width: 100%;
    }
</style>
