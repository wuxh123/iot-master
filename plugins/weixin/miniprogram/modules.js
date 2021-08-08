module.exports = {
    //数据分析
    analysis: {
      getMonthlyRetain: "/datacube/getweanalysisappidmonthlyretaininfo",
      getWeeklyRetain: "/datacube/getweanalysisappidweeklyretaininfo",
      getDailyRetain: "/datacube/getweanalysisappiddailyretaininfo",
      getMonthlyVisitTrend: "/datacube/getweanalysisappidmonthlyvisittrend",
      getWeeklyVisitTrend: "/datacube/getweanalysisappidweeklyvisittrend",
      getDailyVisitTrend: "/datacube/getweanalysisappiddailyvisittrend",
      getUserPortrait: "/datacube/getweanalysisappiduserportrait",
      getVisitDistribution: "/datacube/getweanalysisappidvisitdistribution",
      getVisitPage: "/datacube/getweanalysisappidvisitpage",
      getDailySummary: "/datacube/getweanalysisappiddailysummarytrend"
    },
  
    //客服消息
    customerServiceMessage: {
      setTyping: "/cgi-bin/message/custom/typing",
      uploadTempMedia: "/cgi-bin/media/upload",
      getTempMedia: "/cgi-bin/media/get",
      send: "/cgi-bin/message/custom/send"
    },
  
    //模板消息（已废弃）
    templateMessage: {
      addTemplate: "/cgi-bin/wxopen/template/add",
      deleteTemplate: "/cgi-bin/wxopen/template/del",
      getTemplateLibraryById: "/cgi-bin/wxopen/template/library/get",
      getTemplateLibraryList: "/cgi-bin/wxopen/template/library/list",
      getTemplateList: "/cgi-bin/wxopen/template/list",
      send: "/cgi-bin/message/wxopen/template/send"
    },
  
    //统一服务消息
    uniformMessage: {
      send: "/cgi-bin/message/wxopen/template/uniform_send"
    },
  
    //动态消息
    updatableMessage: {
      createActivityId: "/cgi-bin/message/wxopen/activityid/create",
      setUpdatableMsg: "/cgi-bin/message/wxopen/updatablemsg/send"
    },
  
    //插件管理
    pluginManager: {
      applyPlugin: "/wxa/plugin",
      getPluginDevApplyList: "/wxa/devplugin",
      getPluginList: "/wxa/plugin",
      setDevPluginApplyStatus: "/wxa/devplugin",
      unbindPlugin: "/wxa/plugin"
    },
  
    //附近小程序
    nearByPoi: {
      add: "/wxa/addnearbypoi",
      delete: "/wxa/delnearbypoi",
      getList: "/wxa/getnearbypoilist",
      setShowStatus: "/wxa/setnearbypoishowstatus"
    },
  
  
    //小程序码
    wxacode: {
      createQRCode: "/cgi-bin/wxaapp/createwxaqrcode",
      get: "/wxa/getwxacode",
      getUnlimited: "/wxa/getwxacodeunlimit"
    },
  
    //内容安全
    security: {
      imgSecCheck: "/wxa/img_sec_check",
      mediaCheckAsync: "/wxa/media_check_async",
      msgSecCheck: "/wxa/msg_sec_check"
    },
  
    //生物认证
    soter: {
      verifySignature: "/cgi-bin/soter/verify_signature"
    },
  
    //订阅消息（已废弃）
    subscribeMessage: {
      addTemplate: "/wxaapi/newtmpl/addtemplate",
      deleteTemplate: "/wxaapi/newtmpl/deltemplate",
      getCategory: "/wxaapi/newtmpl/getcategory",
      getPubTemplateKeyWordsById: "/wxaapi/newtmpl/getpubtemplatekeywords",
      getPubTemplateTitleList: "/wxaapi/newtmpl/getpubtemplatetitles",
      getTemplateList: "/wxaapi/newtmpl/gettemplate",
      send: "/cgi-bin/message/subscribe/send"
    },
  
  };