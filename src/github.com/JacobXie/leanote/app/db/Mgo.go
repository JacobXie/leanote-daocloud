package db

import (
	"fmt"
	. "github.com/JacobXie/leanote/app/lea"
	"github.com/JacobXie/leanote/app/info"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

// Init mgo and the common DAO

// 数据连接
var Session *mgo.Session

// 各个表的Collection对象
var Notebooks *mgo.Collection
var Notes *mgo.Collection
var NoteContents *mgo.Collection
var NoteContentHistories *mgo.Collection

var ShareNotes *mgo.Collection
var ShareNotebooks *mgo.Collection
var HasShareNotes *mgo.Collection
var Blogs *mgo.Collection
var Users *mgo.Collection
var Groups *mgo.Collection
var GroupUsers *mgo.Collection

var Tags *mgo.Collection
var NoteTags *mgo.Collection
var TagCounts *mgo.Collection

var UserBlogs *mgo.Collection

var Tokens *mgo.Collection

var Suggestions *mgo.Collection

// Album & file(image)
var Albums *mgo.Collection
var Files *mgo.Collection
var Attachs *mgo.Collection

var NoteImages *mgo.Collection
var Configs *mgo.Collection
var EmailLogs *mgo.Collection

// blog
var BlogLikes *mgo.Collection
var BlogComments *mgo.Collection
var Reports *mgo.Collection
var BlogSingles *mgo.Collection
var Themes *mgo.Collection

// session
var Sessions *mgo.Collection

// 初始化时连接数据库
func Init(url, dbname string) {
	ok := true
	config := revel.Config
	if url == "" {
		url, ok = config.String("db.url")
		if !ok {
			url, ok = config.String("db.urlEnv")
			if ok {
				Log("get db conf from urlEnv: " + url)
			}
		} else {
			Log("get db conf from db.url: " + url)
		}

		if ok {
			// get dbname from urlEnv
			urls := strings.Split(url, "/")
			dbname = urls[len(urls)-1]
		}
	}
	if dbname == "" {
		dbname, _ = config.String("db.dbname")
		Log("dbname : " + dbname)
	}

	// get db config from host, port, username, password
	if !ok {
		host, _ := revel.Config.String("db.host")
		port, _ := revel.Config.String("db.port")
		username, _ := revel.Config.String("db.username")
		password, _ := revel.Config.String("db.password")
		usernameAndPassword := username + ":" + password + "@"
		if username == "" || password == "" {
			usernameAndPassword = ""
		}
		//mpodified by JacobXie
		//url = "mongodb://" + usernameAndPassword + host + ":" + port + "/" + dbname
    url = fmt.Sprintf("%s%s:%s/%s", usernameAndPassword,host, port, dbname)
	}
	Log(url)

	// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	// mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
	var err error
	Session, err = mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	Session.SetMode(mgo.Monotonic, true)

	// notebook
	Notebooks = Session.DB(dbname).C("notebooks")

	// notes
	Notes = Session.DB(dbname).C("notes")

	// noteContents
	NoteContents = Session.DB(dbname).C("note_contents")
	NoteContentHistories = Session.DB(dbname).C("note_content_histories")

	// share
	ShareNotes = Session.DB(dbname).C("share_notes")
	ShareNotebooks = Session.DB(dbname).C("share_notebooks")
	HasShareNotes = Session.DB(dbname).C("has_share_notes")

	// user
	Users = Session.DB(dbname).C("users")
	// group
	Groups = Session.DB(dbname).C("groups")
	GroupUsers = Session.DB(dbname).C("group_users")

	// blog
	Blogs = Session.DB(dbname).C("blogs")

	// tag
	Tags = Session.DB(dbname).C("tags")
	NoteTags = Session.DB(dbname).C("note_tags")
	TagCounts = Session.DB(dbname).C("tag_count")

	// blog
	UserBlogs = Session.DB(dbname).C("user_blogs")
	BlogSingles = Session.DB(dbname).C("blog_singles")
	Themes = Session.DB(dbname).C("themes")

	// find password
	Tokens = Session.DB(dbname).C("tokens")

	// Suggestion
	Suggestions = Session.DB(dbname).C("suggestions")

	// Album & file
	Albums = Session.DB(dbname).C("albums")
	Files = Session.DB(dbname).C("files")
	Attachs = Session.DB(dbname).C("attachs")

	NoteImages = Session.DB(dbname).C("note_images")

	Configs = Session.DB(dbname).C("configs")
	EmailLogs = Session.DB(dbname).C("email_logs")

	// 社交
	BlogLikes = Session.DB(dbname).C("blog_likes")
	BlogComments = Session.DB(dbname).C("blog_comments")

	// 举报
	Reports = Session.DB(dbname).C("reports")

	// session
	Sessions = Session.DB(dbname).C("sessions")

	//Modified by JacobXie
	//系统初始化，即初始化管理员和一些基本数据
	countNum, err := Users.Count()
	if err != nil {
		panic(err)
	}
	if countNum == 0 {
		//初始化管理员，admin，密码 admin123
		user := info.User{}
		user.UserId = bson.NewObjectId()
		user.Email = "admin@leanote.com"
		user.Verified = false
		user.Username,ok = revel.Config.String("adminUsername")
		if !ok || user.Username == ""{
			user.Username = "admin"
		}
		user.UsernameRaw = user.Username
		user.Pwd = GenPwd(user.Username+"123")
		user.CreatedTime = time.Now()
		user.Theme = "simple"
		user.NotebookWidth = 160
		user.NoteListWidth = 266
		Insert(Users, user)

		//Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"openRegister", ValueStr:"open" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"toImageBinPath", ValueStr:"lllllllllll" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"noteSubDomain", ValueStr:"" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"blogSubDomain", ValueStr:"" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"leaSubDomain", ValueStr:"" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"recommendTags", ValueArr:[]string{"小写", "golang", "leanote",} ,IsArr:true})
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"newTags", ValueArr:[]string{"小写", "golang", "leanote", "haha",} , IsArr:true })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailHost", ValueStr:"smtp.163.com" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailPort", ValueStr:"25" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailUsername", ValueStr:"" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailPassword", ValueStr:"" })

		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateHeader", ValueStr: "<div style=\"width: 600px; margin:auto; border-radius:5px; border: 1px solid #ccc; padding: 20px;\">\r\n\t\t\t<div>\r\n\t\t\t\t<div>\r\n\t\t\t\t\t<div style=\"float:left; height: 40px;\">\r\n\t\t\t\t\t\t<a href=\"{{$.siteUrl}}\" style=\"font-size: 24px\">leanote</a>\r\n\t\t\t\t\t</div>\r\n\t\t\t\t\t<div style=\"float:left; height:40px; line-height:40px;\">\r\n\t\t\t\t\t\t&nbsp;&nbsp;| &nbsp;<span style=\"font-size:14px\">{{$.subject}}</span>\r\n\t\t\t\t\t</div>\r\n\t\t\t\t\t<div style=\"clear:both\"></div>\r\n\t\t\t\t</div>\r\n\t\t\t</div>\r\n\t\t\t<hr style=\"border:none;border-top: 1px solid #ccc\"/>\r\n\t\t\t<div style=\"margin-top: 20px; font-size: 14px;\">\r\n\t\t\t\t" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateFooter", ValueStr:"</div>\r\n\r\n\t\t\t<div id=\"leanoteFooter\" style=\"margin-top: 30px; border-top: 1px solid #ccc\">\r\n\t\t\t\t<style>\r\n\t\t\t\t\t#leanoteFooter {\r\n\t\t\t\t\t\tcolor: #666;\r\n\t\t\t\t\t\tfont-size: 12px;\r\n\t\t\t\t\t}\r\n\t\t\t\t\t#leanoteFooter a {\r\n\t\t\t\t\t\tcolor: #666;\r\n\t\t\t\t\t\tfont-size: 12px;\r\n\t\t\t\t\t}\r\n\t\t\t\t</style>\r\n\t\t\t\t<a href=\"{{$.siteUrl}}\">leanote</a>, your own cloud note!\r\n\t\t\t</div>\r\n\t\t</div>" })

		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateRegisterSubject", ValueStr:"欢迎来到leanote, 请验证邮箱" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateRegister", ValueStr:"{{header}}\r\n<p>\r\n{{$.user.email}} 您好, 欢迎来到leanote. \r\n</p>\r\n<p>\r\n请点击链接验证邮箱: <a href=\"{{$.tokenUrl}}\">{{$.tokenUrl}}</a>\r\n</p>\r\n<p>\r\n{{$.tokenTimeout}}小时后过期.\r\n</p>\r\n{{footer}}"})

		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateFindPasswordSubject", ValueStr:"找回密码" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateFindPassword", ValueStr:"{{header}}\r\n<p>\r\n请点击链接修改密码 <a href=\"{{$.tokenUrl}}\">{{$.tokenUrl}}</a>\r\n</p>\r\n<p>\r\n{{$.tokenTimeout}}小时后过期.\r\n</p>\r\n\r\n{{footer}}"})

		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateUpdateEmailSubject", ValueStr:"验证邮箱" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateUpdateEmail", ValueStr:"{{header}}\r\n<p>\r\n邮箱验证后您的登录邮箱为: {{$.newEmail}}\r\n</p>\r\n<p>\r\n请点击链接验证邮箱: <a href=\"{{$.tokenUrl}}\">{{$.tokenUrl}}</a>\r\n</p>\r\n<p>\r\n{{$.tokenTimeout}}小时后过期.\r\n</p>\r\n{{footer}}\r\n"})

		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateInviteSubject", ValueStr:"邀请注册leanote" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateInvite", ValueStr:"{{header}}\r\n\r\n<p>您好, 您的好友{{$.user.email}}邀请您注册leanote</p>\r\n\r\n<p>Ta的留言: {{$.content}}</p>\r\n\r\n<p>点击链接注册leanote <a href=\"{{$.registerUrl}}\">{{$.registerUrl}}</a></p>\r\n\r\n{{footer}}\r\n"})

		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateCommentSubject", ValueStr:"评论提醒" })
		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"emailTemplateComment", ValueStr:"{{header}}\r\n<p>\r\n{{if $.commentedUser.isBlogAuthor}}\r\n您的博客 \"{{$.blog.title}}\" 被 {{$.commentUser.username}} 评论了.\r\n{{else}}\r\n您在 \"{{$.blog.title}}\" 发表的评论被 {{$.commentUser.username}}{{if $.commentUser.isBlogAuthor}}(作者){{end}} 评论了.\r\n{{end}}\r\n</p>\r\n\r\n<div>\r\n<b>评论内容: </b>\r\n<blockquote>{{$.commentContent}}</blockquote>\r\n</div>\r\n<p>\r\n博客链接: <a href=\"{{$.blog.url}}\">{{$.blog.url}}</a>\r\n</p>\r\n{{footer}} "})

		Insert(Configs,& info.Config{ConfigId:bson.NewObjectId(), UserId:user.UserId, Key:"userFilterEmail", ValueStr:"" })


		theme_elegant := info.Theme{}
		theme_elegant.ThemeId = bson.NewObjectId()
		theme_elegant.UserId = user.UserId
		theme_elegant.Name = "leanote elegant"
		theme_elegant.Version = "1.0"
		theme_elegant.Author = "leanote.com"
		theme_elegant.AuthorUrl = "http://leanote.com"
		theme_elegant.Path =  "public/blog/themes/elegant"
		theme_elegant.Info =  map[string]interface{}{ "Version" : "1.0",
														"Author" : "leanote.com",
														"AuthorUrl" : "http://leanote.com",
														"FriendLinks" : []map[string]string{
																			//map[string]string{ "Title" : "我的笔记", "Url" : "http://leanote.com/note" },
																			//map[string]string{ "Title" : "leanote home", "Url" : "http://leanote.com" },
																			//map[string]string{ "Title" : "leanote 社区", "Url" : "http://bbs.leanote.com" },
																			map[string]string{ "Title" : "lea++", "Url" : "http://lea.leanote.com" },
																			//map[string]string{ "Title" : "leanote github", "Url" : "https://github.com/leanote/leanote" },
																		},
														"Name" : "leanote elegant",
														}
		theme_elegant.IsActive =  true
		theme_elegant.IsDefault =  true
		theme_elegant.Style = "blog_daqi"
		theme_elegant.CreatedTime = time.Now()
		theme_elegant.UpdatedTime = theme_elegant.CreatedTime
		Insert(Themes,theme_elegant)

		theme_default := info.Theme{}
		theme_default.ThemeId = bson.NewObjectId()
		theme_default.UserId = user.UserId
		theme_default.Name = "leanote default theme"
		theme_default.Version = "1.0"
		theme_default.Author =  "leanote.com"
		theme_default.AuthorUrl = "http://leanote.com"
		theme_default.Path = "public/blog/themes/default"

		theme_default.Info = map[string]interface{}{ "AuthorUrl" : "http://leanote.com",
														"FriendLinks" : []map[string]string{
																			//map[string]string{ "Url" : "http://leanote.com/note", "Title" : "我的笔记" },
																			//map[string]string{ "Title" : "leanote home", "Url" : "http://leanote.com" },
																			//map[string]string{ "Title" : "leanote 社区", "Url" : "http://bbs.leanote.com" },
																			map[string]string{ "Url" : "http://lea.leanote.com", "Title" : "lea++" },
																			//map[string]string{ "Title" : "leanote github", "Url" : "https://github.com/leanote/leanote" },
																		},
														"Name" : "leanote default theme",
														"Version" : "1.0",
														"Author" : "leanote.com",
													}
		theme_default.IsActive = false
		theme_default.IsDefault =  true
		theme_default.Style = "blog_default"
		theme_default.CreatedTime = time.Now()
		theme_default.UpdatedTime = theme_default.CreatedTime
		Insert(Themes,theme_default)

		theme_nav_fixed := info.Theme{}
		theme_nav_fixed.ThemeId = bson.NewObjectId()
		theme_nav_fixed.UserId = user.UserId
		theme_nav_fixed.Name = "leanote nav fixed"
		theme_nav_fixed.Version  = "1.0"
		theme_nav_fixed.Author = "leanote.com"
		theme_nav_fixed.AuthorUrl = "http://leanote.com"
		theme_nav_fixed.Path = "public/blog/themes/nav_fixed"
		theme_nav_fixed.Info = map[string]interface{}{
														"Name" : "leanote nav fixed",
														"Version" : "1.0",
														"Author" : "leanote.com",
														"AuthorUrl" : "http://leanote.com",
														"FriendLinks" :
																		[]map[string]string{
																			//map[string]string{ "Title" : "我的笔记","Url" : "http://leanote.com/note" },
																			//map[string]string{ "Title" : "leanote home", "Url" : "http://leanote.com" },
																			//map[string]string{ "Title" : "leanote 社区", "Url" : "http://bbs.leanote.com" },
																			map[string]string{ "Title" : "lea++", "Url" : "http://lea.leanote.com" },
																			//map[string]string{ "Title" : "leanote github", "Url" : "https://github.com/leanote/leanote" },
																	},
														}
		theme_nav_fixed.IsActive = false
		theme_nav_fixed.IsDefault = true
		theme_nav_fixed.Style = "blog_left_fixed"
		theme_nav_fixed.CreatedTime = time.Now()
		theme_nav_fixed.UpdatedTime = theme_nav_fixed.CreatedTime
		Insert(Themes,theme_nav_fixed)


		blog_single := info.BlogSingle{}
		blog_single.SingleId = bson.NewObjectId()
		blog_single.UserId = user.UserId
		blog_single.Title = "About Me"
		blog_single.UrlTitle = "About-Me"
		blog_single.Content = "<p>Hello,&nbsp;I am Leanote (^_^).</p>"
		blog_single.CreatedTime = time.Now()
		blog_single.UpdatedTime = blog_single.CreatedTime
		Insert(BlogSingles,blog_single)

		Insert(Notebooks, & info.Notebook{NotebookId:bson.NewObjectId(), UserId:user.UserId, Seq:-1, Title:"Life", UrlTitle:"life", CreatedTime:time.Now(), UpdatedTime: time.Now()})
		Insert(Notebooks, & info.Notebook{NotebookId:bson.NewObjectId(), UserId:user.UserId, Seq:-1, Title:"Work", UrlTitle:"work", CreatedTime:time.Now(), UpdatedTime: time.Now()})
		Insert(Notebooks, & info.Notebook{NotebookId:bson.NewObjectId(), UserId:user.UserId, Seq:-1, Title:"Study", UrlTitle:"study", CreatedTime:time.Now(), UpdatedTime: time.Now()})
		Insert(Notebooks, & info.Notebook{NotebookId:bson.NewObjectId(), UserId:user.UserId, Seq:-1, Title:"Other", UrlTitle:"other", CreatedTime:time.Now(), UpdatedTime: time.Now()})

		tag := info.Tag{}
		tag.UserId = user.UserId
		Insert(Tags,tag)

		token := info.Token{}
		token.UserId = user.UserId
		token.Email = "admin@leanote.com"
		token.Token = NewGuidWith(token.Email)
		token.Type = 1
		token.CreatedTime =  user.CreatedTime
		Insert(Tokens,token)


		user_blog := info.UserBlog{}
		user_blog.UserId = user.UserId
		user_blog.Title = "Leanote's Blog"
		user_blog.SubTitle = "I love Leanote!"
		user_blog.AboutMe = "<p>Hello, 大家好, 我是leanote, 赶紧来体验leanote吧!!</p>"
		user_blog.CanComment = true
		user_blog.CommentType = "default"
		user_blog.DisqusId =  "leanote"
		user_blog.Style = "blog_daqi"
		user_blog.ThemeId = theme_elegant.ThemeId
		user_blog.Singles = []map[string]string{
													map[string]string{ "SingleId" : blog_single.SingleId.Hex(), "Title" : "About Me", "UrlTitle" : "About-Me" },
												}
		Insert(UserBlogs,user_blog)

	}
}

func close() {
	Session.Close()
}

// common DAO
// 公用方法

//----------------------

func Insert(collection *mgo.Collection, i interface{}) bool {
	err := collection.Insert(i)
	return Err(err)
}

//----------------------

// 适合一条记录全部更新
func Update(collection *mgo.Collection, query interface{}, i interface{}) bool {
	err := collection.Update(query, i)
	return Err(err)
}
func Upsert(collection *mgo.Collection, query interface{}, i interface{}) bool {
	_, err := collection.Upsert(query, i)
	return Err(err)
}
func UpdateAll(collection *mgo.Collection, query interface{}, i interface{}) bool {
	_, err := collection.UpdateAll(query, i)
	return Err(err)
}
func UpdateByIdAndUserId(collection *mgo.Collection, id, userId string, i interface{}) bool {
	err := collection.Update(GetIdAndUserIdQ(id, userId), i)
	return Err(err)
}

func UpdateByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) bool {
	err := collection.Update(GetIdAndUserIdBsonQ(id, userId), i)
	return Err(err)
}
func UpdateByIdAndUserIdField(collection *mgo.Collection, id, userId, field string, value interface{}) bool {
	return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}
func UpdateByIdAndUserIdMap(collection *mgo.Collection, id, userId string, v bson.M) bool {
	return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": v})
}

func UpdateByIdAndUserIdField2(collection *mgo.Collection, id, userId bson.ObjectId, field string, value interface{}) bool {
	return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}
func UpdateByIdAndUserIdMap2(collection *mgo.Collection, id, userId bson.ObjectId, v bson.M) bool {
	return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": v})
}

//
func UpdateByQField(collection *mgo.Collection, q interface{}, field string, value interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": bson.M{field: value}})
	return Err(err)
}
func UpdateByQI(collection *mgo.Collection, q interface{}, v interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": v})
	return Err(err)
}

// 查询条件和值
func UpdateByQMap(collection *mgo.Collection, q interface{}, v interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": v})
	return Err(err)
}

//------------------------

// 删除一条
func Delete(collection *mgo.Collection, q interface{}) bool {
	err := collection.Remove(q)
	return Err(err)
}
func DeleteByIdAndUserId(collection *mgo.Collection, id, userId string) bool {
	err := collection.Remove(GetIdAndUserIdQ(id, userId))
	return Err(err)
}
func DeleteByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId) bool {
	err := collection.Remove(GetIdAndUserIdBsonQ(id, userId))
	return Err(err)
}

// 删除所有
func DeleteAllByIdAndUserId(collection *mgo.Collection, id, userId string) bool {
	_, err := collection.RemoveAll(GetIdAndUserIdQ(id, userId))
	return Err(err)
}
func DeleteAllByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId) bool {
	_, err := collection.RemoveAll(GetIdAndUserIdBsonQ(id, userId))
	return Err(err)
}

func DeleteAll(collection *mgo.Collection, q interface{}) bool {
	_, err := collection.RemoveAll(q)
	return Err(err)
}

//-------------------------

func Get(collection *mgo.Collection, id string, i interface{}) {
	collection.FindId(bson.ObjectIdHex(id)).One(i)
}
func Get2(collection *mgo.Collection, id bson.ObjectId, i interface{}) {
	collection.FindId(id).One(i)
}

func GetByQ(collection *mgo.Collection, q interface{}, i interface{}) {
	collection.Find(q).One(i)
}
func ListByQ(collection *mgo.Collection, q interface{}, i interface{}) {
	collection.Find(q).All(i)
}

func ListByQLimit(collection *mgo.Collection, q interface{}, i interface{}, limit int) {
	collection.Find(q).Limit(limit).All(i)
}

// 查询某些字段, q是查询条件, fields是字段名列表
func GetByQWithFields(collection *mgo.Collection, q bson.M, fields []string, i interface{}) {
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	collection.Find(q).Select(selector).One(i)
}

// 查询某些字段, q是查询条件, fields是字段名列表
func ListByQWithFields(collection *mgo.Collection, q bson.M, fields []string, i interface{}) {
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	collection.Find(q).Select(selector).All(i)
}
func GetByIdAndUserId(collection *mgo.Collection, id, userId string, i interface{}) {
	collection.Find(GetIdAndUserIdQ(id, userId)).One(i)
}
func GetByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) {
	collection.Find(GetIdAndUserIdBsonQ(id, userId)).One(i)
}

// 按field去重
func Distinct(collection *mgo.Collection, q bson.M, field string, i interface{}) {
	collection.Find(q).Distinct(field, i)
}

//----------------------

func Count(collection *mgo.Collection, q interface{}) int {
	cnt, err := collection.Find(q).Count()
	if err != nil {
		Err(err)
	}
	return cnt
}

func Has(collection *mgo.Collection, q interface{}) bool {
	if Count(collection, q) > 0 {
		return true
	}
	return false
}

//-----------------

// 得到主键和userId的复合查询条件
func GetIdAndUserIdQ(id, userId string) bson.M {
	return bson.M{"_id": bson.ObjectIdHex(id), "UserId": bson.ObjectIdHex(userId)}
}
func GetIdAndUserIdBsonQ(id, userId bson.ObjectId) bson.M {
	return bson.M{"_id": id, "UserId": userId}
}

// DB处理错误
func Err(err error) bool {
	if err != nil {
		fmt.Println(err)
		// 删除时, 查找
		if err.Error() == "not found" {
			return true
		}
		return false
	}
	return true
}

// 检查mognodb是否lost connection
// 每个请求之前都要检查!!
func CheckMongoSessionLost() {
	// fmt.Println("检查CheckMongoSessionLostErr")
	err := Session.Ping()
	if err != nil {
		Log("Lost connection to db!")
		Session.Refresh()
		err = Session.Ping()
		if err == nil {
			Log("Reconnect to db successful.")
		} else {
			Log("重连失败!!!! 警告")
		}
	}
}
