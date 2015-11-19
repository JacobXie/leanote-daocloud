var MSG={"app":"Leanote","share":"分享","noTag":"无标签","inputEmail":"请输入Email","history":"历史记录","editorTips":"帮助","editorTipsInfo":"<h4>1. 快捷键</h4>ctrl+shift+c 代码块切换 <h4>2. shift+enter 跳出当前区域</h4>比如在代码块中<img src=\"/images/outofcode.png\" style=\"width: 90px\"/>按shift+enter可跳出当前代码块.","all":"最新","trash":"废纸篓","delete":"删除","unTitled":"无标题","defaultShare":"默认分享","writingMode":"写作模式","normalMode":"普通模式","saving":"正在保存","saveSuccess":"保存成功","Are you sure to delete it ?":"确认删除?","Insert link into content":"将附件链接插入到内容中","Download":"下载","Delete":"删除","update":"更新","Update Time":"更新时间","Create Time":"创建时间","Post Url":"博文链接","close":"关闭","cancel":"取消","send":"发送","shareToFriends":"分享给好友","publicAsBlog":"公开为博客","cancelPublic":"取消公开为博客","move":"移动","copy":"复制","rename":"重命名","exportPdf":"导出PDF","addChildNotebook":"添加子笔记本","deleteAllShared":"删除所有共享","deleteSharedNotebook":"删除共享笔记本","copyToMyNotebook":"复制到我的笔记本","sendSuccess":"发送成功","friendEmail":"好友邮箱","readOnly":"只读","writable":"可写","inputFriendEmail":"请输入好友邮箱","clickToChangePermission":"点击改变权限","sendInviteEmailToYourFriend":"发送邀请email给Ta","friendNotExits":"该用户还没有注册%s, 复制邀请链接发送给Ta, 邀请链接: %s","emailBodyRequired":"邮件内容不能为空","inviteEmailBody":"Hi, 你好, 我是%s, %s非常好用, 快来注册吧!","historiesNum":"leanote会保存笔记的最近<b>10</b>份历史记录","noHistories":"无历史记录","datetime":"日期","restoreFromThisVersion":"从该版本还原","confirmBackup":"确定要从该版还原? 还原前leanote会备份当前版本到历史记录中.","errorEmail":"请输入正确的email","Hyperlink":"超链接","Please provide the link URL and an optional title":"请填写链接和一个可选的标题","optional title":"可选标题","Cancel":"取消","Strong":"粗体","strong text":"粗体","Emphasis":"斜体","emphasized text":"斜体","Blockquote":"引用","Code Sample":"代码","enter code here":"代码","Image":"图片","Heading":"标题","Numbered List":"有序列表","Bulleted List":"无序列表","List item":"项目","Horizontal Rule":"水平线","Undo":"撤销","Redo":"重做","enter image description here":"图片标题","enter link description here":"链接标题","Add Album":"添加相册","Cannot delete default album":"不能删除默认相册","Cannot rename default album":"不能重命名默认相册","Rename Album":"重命名","Add Success!":"添加成功!","Rename Success!":"重命名成功!","Delete Success!":"删除成功","Are you sure to delete this image ?":"确定删除该图片?","click to remove this image":"删除图片","error":"错误","Error":"错误","Prev":"上一页","Next":"下一页"};function getMsg(key, data) {var msg = MSG[key];if(msg) {if(data) {if(!isArray(data)) {data = [data];}for(var i = 0; i < data.length; ++i) {msg = msg.replace("%s", data[i]);}}return msg;}return key;}