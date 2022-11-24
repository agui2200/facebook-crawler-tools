package fb

import (
	"encoding/base64"
	"facebook_login/global"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/google/uuid"
	"github.com/imroc/req"
)

type LoginParam struct {
	Datr      string `json:"datr"`
	Lsd       string `json:"lsd"`
	KeyID     string `json:"keyId"`
	PublicKey string `json:"publicKey"`
	Jazoest   string `json:"jazoest"`

	Cookie string `json:"cookie"`
}

type GraphqlParam struct {
	UserId    string
	Dyn       string
	Csr       string
	Rev       string
	Hsi       string
	Comet_req string
	Fb_dtsg   string
	Jazoes    string
	Lsd       string
	BaseDyn   []int
	BaseCsr   []int

	//需手动设置参数
	Cookie                   string
	Doc_id                   string
	Fb_api_req_friendly_name string
	Fb_api_caller_class      string
	Variables                string
	Ccg                      string
}

type GroupsInfo struct {
	Id     string
	Name   string
	Url    string
	Cursor string
}

func Init(proxy string) {
	if len(proxy) > 0 {
		//初始化一些Http设置
		req.SetProxy(func(r *http.Request) (*url.URL, error) {
			if strings.Contains(r.URL.Hostname(), "facebook") {
				return url.Parse(proxy)
			}
			return nil, nil
		})
	}
	req.EnableCookie(false)
	req.SetTimeout(60 * time.Second)
	req.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

func GetGraphqlParam(cookie string) (bool, GraphqlParam) {
	var data GraphqlParam
	header := req.Header{
		`accept`:                    ` text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		`accept-language`:           ` zh-CN,zh;q=0.9`,
		`cache-control`:             ` no-cache`,
		`pragma`:                    ` no-cache`,
		`sec-ch-ua`:                 ` "Chromium";v="104", " Not A;Brand";v="99", "Google Chrome";v="104"`,
		`sec-ch-ua-mobile`:          ` ?0`,
		`sec-ch-ua-platform`:        ` "Windows"`,
		`sec-fetch-dest`:            ` document`,
		`sec-fetch-mode`:            ` navigate`,
		`sec-fetch-site`:            ` none`,
		`sec-fetch-user`:            ` ?1`,
		`upgrade-insecure-requests`: ` 1`,
		`cookie`:                    cookie,
		`user-agent`:                ` Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`,
	}

	res, err1 := req.Get("https://www.facebook.com/", header)
	if err1 != nil {
		return false, GraphqlParam{}
	}

	buffer, _ := io.ReadAll(res.Response().Body)
	body := string(buffer)

	//准备参数
	reg := regexp.MustCompile(`<script id="__eqmc" type="application/json" nonce=".*?">(.*?)</script>`)
	result := reg.FindAllStringSubmatch(body, -1)

	if strings.Contains(body, "Sorry") || len(result) <= 0 {
		return false, GraphqlParam{}
	}

	userId := global.StrBetween(cookie, `c_user=`, `;`)
	baseDyn, dyn := getDyn(body)
	baseCsr, csr := getCsr(body)
	rev := global.StrBetween(body, `data-btmanifest="`, `"`)
	rev = strings.ReplaceAll(rev, "_main", "")
	comet_req := global.StrBetween(result[0][1], `comet_req=`, `&`)
	hsi := global.StrBetween(result[0][1], `"e":"`, `",`)
	fb_dtsg := global.StrBetween(result[0][1], `"f":"`, `",`)
	jazoest := global.StrBetween(result[0][1], `jazoest=`, `",`)
	lsd := global.StrBetween(body, `"LSD",[],{"token":"`, `"},`)

	data.UserId = userId
	data.Dyn = dyn
	data.Csr = csr
	data.Rev = rev
	data.Comet_req = comet_req
	data.Hsi = hsi
	data.Fb_dtsg = fb_dtsg
	data.Jazoes = jazoest
	data.Lsd = lsd
	data.BaseCsr = baseCsr
	data.BaseDyn = baseDyn
	data.Cookie = cookie

	return true, data
}

func Login(user string, pwd string) (bool, string) {

	var param = getLoginParam()
	if len(param.Datr) == 0 || len(param.KeyID) == 0 || len(param.Lsd) == 0 || len(param.PublicKey) == 0 || len(param.KeyID) == 0 {
		return false, "GetLoginParam is empty!"
	}

	var timeStamp = time.Now().Unix()
	var publicKey = param.PublicKey
	var keyId = param.KeyID

	var encrypted = encpass(pwd, strconv.FormatInt(timeStamp, 10), publicKey, keyId)

	if len(encrypted) == 0 {
		return false, "Encpass error!"
	}

	var creation_time_base64 = base64.StdEncoding.EncodeToString([]byte(`{"type":0,"creation_time":` + strconv.FormatInt(timeStamp-4, 10) + `,"callsite_id":381229079575946}`))
	var header = req.Header{
		`accept`:                    ` text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		`accept-language`:           ` zh-CN,zh;q=0.9`,
		`cache-control`:             ` no-cache`,
		`content-type`:              ` application/x-www-form-urlencoded`,
		`cookie`:                    param.Cookie + `wd=1253x937;`,
		`origin`:                    ` https://www.facebook.com`,
		`pragma`:                    ` no-cache`,
		`referer`:                   ` https://www.facebook.com/`,
		`sec-ch-ua`:                 ` "Chromium";v="104", " Not A;Brand";v="99", "Google Chrome";v="104"`,
		`sec-ch-ua-mobile`:          ` ?0`,
		`sec-ch-ua-platform`:        ` "Windows"`,
		`sec-fetch-dest`:            ` document`,
		`sec-fetch-mode`:            ` navigate`,
		`sec-fetch-site`:            ` same-origin`,
		`sec-fetch-user`:            ` ?1`,
		`upgrade-insecure-requests`: ` 1`,
		`user-agent`:                ` Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`,
	}

	var data = req.Param{
		`jazoest`:      param.Jazoest,
		`lsd`:          param.Lsd,
		`email`:        user,
		`login_source`: ` comet_headerless_login`,
		`next`:         ` `,
		`encpass`:      encrypted,
	}

	res, err := req.Post("https://www.facebook.com/login/?privacy_mutation_token="+creation_time_base64,
		header, data)
	if err != nil {
		return false, "Http request error"
	}

	buffer, _ := io.ReadAll(res.Response().Body)
	var body = string(buffer)

	var cookie string
	for _, item := range res.Response().Cookies() {
		cookie += item.Name + "=" + item.Value + ";"
	}

	if strings.Contains(body, "你暂时被禁止使用此功能") {
		return false, "你暂时被禁止使用此功能"
	}
	if strings.Contains(body, "你目前没有访问公共主页的权限") {
		return false, "你目前没有访问公共主页的权限"
	}
	if strings.Contains(body, "帐号或密码无效") {
		return false, "帐号或密码无效"
	}
	if strings.Contains(body, "你输入的邮箱或手机号未绑定任何帐户") {
		return false, "你输入的邮箱或手机号未绑定任何帐户"
	}
	if strings.Contains(body, "无法处理你的请求") {
		return false, "无法处理你的请求"
	}
	if res.Response().StatusCode == 302 && strings.Contains(cookie, "c_user") {
		return true, cookie
	}

	log.Println(body)
	return false, "未知错误"
}

func Like(param GraphqlParam, doc_id string, feedback_id string, encrypted_tracking string) (bool, string) {

	//1678524932434102  大爱 1635855486666999  普通赞
	likeType := "1635855486666999"
	milliTimestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	guid := uuid.New().String()
	variables := `{"input":{"attribution_id_v2":"CometHomeRoot.react,comet.home,via_cold_start,` + milliTimestamp + `,420759,4748854339,","feedback_id":"` + feedback_id + `","feedback_reaction_id":"` + likeType + `","feedback_source":"NEWS_FEED","is_tracking_encrypted":true,"tracking":["` + encrypted_tracking + `"],"session_id":"` + guid + `","actor_id":"` + param.UserId + `","client_mutation_id":"` + global.RandomStr(1) + `"},"useDefaultActor":false,"scale":1.5}`

	param.Fb_api_req_friendly_name = "CometUFIFeedbackReactMutation"
	param.Variables = variables
	param.Fb_api_caller_class = "RelayModern"
	param.Ccg = "GOOD"
	param.Doc_id = doc_id

	return graphql(param)
}

func Post(param GraphqlParam, doc_id string, text string, base_state int) (bool, string) {

	var state = "SELF"
	if base_state == 0 {
		//帖子所有人可见
		state = "EVERYONE"
	}
	if base_state == 1 {
		//帖子仅自己可见
		state = "SELF"
	}

	milliTimestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	guid := uuid.New().String()
	variables := `{"input":{"composer_entry_point":"inline_composer","composer_source_surface":"newsfeed","composer_type":"feed","idempotence_token":"` + guid + `_FEED","source":"WWW","attachments":[],"audience":{"privacy":{"allow":[],"base_state":"` + state + `","deny":[],"tag_expansion_state":"UNSPECIFIED"}},"message":{"ranges":[],"text":"` + text + `"},"with_tags_ids":[],"inline_activities":[],"explicit_place_id":"0","text_format_preset_id":"0","logging":{"composer_session_id":"` + guid + `"},"navigation_data":{"attribution_id_v2":"CometHomeRoot.react,comet.home,logo,` + milliTimestamp + `,120934,4748854339,"},"tracking":[null],"actor_id":"` + param.UserId + `","client_mutation_id":"` + global.RandomStr(1) + `"},"displayCommentsFeedbackContext":null,"displayCommentsContextEnableComment":null,"displayCommentsContextIsAdPreview":null,"displayCommentsContextIsAggregatedShare":null,"displayCommentsContextIsStorySet":null,"feedLocation":"NEWSFEED","feedbackSource":1,"focusCommentID":null,"gridMediaWidth":null,"groupID":null,"scale":1.5,"privacySelectorRenderLocation":"COMET_STREAM","renderLocation":"homepage_stream","useDefaultActor":false,"inviteShortLinkKey":null,"isFeed":true,"isFundraiser":false,"isFunFactPost":false,"isGroup":false,"isEvent":false,"isTimeline":false,"isSocialLearning":false,"isPageNewsFeed":false,"isProfileReviews":false,"isWorkSharedDraft":false,"UFI2CommentsProvider_commentsKey":"CometModernHomeFeedQuery","hashtag":null,"canUserManageOffers":false,"__relay_internal__pv__FBReelsEnableDeferrelayprovider":false}`

	param.Fb_api_req_friendly_name = "ComposerStoryCreateMutation"
	param.Variables = variables
	param.Fb_api_caller_class = "RelayModern"
	param.Ccg = "EXCELLENT"
	param.Doc_id = doc_id

	success, result := graphql(param)
	if !success {
		return success, result
	}

	//取出发布成功的帖子链接
	json, err := simplejson.NewJson([]byte(result))
	if err != nil {
		return false, "Json parse error"
	}

	url := json.Get("data").Get("story_create").Get("story").Get("url").MustString()
	if len(url) > 0 {
		return true, url
	}

	return false, ""
}

func GetMyGroups(param GraphqlParam, doc_id string, num int) (bool, []GroupsInfo) {

	header := req.Header{
		`accept`:                    ` text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		`accept-language`:           ` zh-CN,zh;q=0.9`,
		`cache-control`:             ` no-cache`,
		`pragma`:                    ` no-cache`,
		`sec-ch-ua`:                 ` "Chromium";v="104", " Not A;Brand";v="99", "Google Chrome";v="104"`,
		`sec-ch-ua-mobile`:          ` ?0`,
		`sec-ch-ua-platform`:        ` "Windows"`,
		`sec-fetch-dest`:            ` document`,
		`sec-fetch-mode`:            ` navigate`,
		`sec-fetch-site`:            ` none`,
		`sec-fetch-user`:            ` ?1`,
		`upgrade-insecure-requests`: ` 1`,
		`cookie`:                    param.Cookie,
		`user-agent`:                ` Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`,
	}

	res, err1 := req.Get("https://www.facebook.com/groups/feed/", header)
	if err1 != nil {
		return false, nil
	}

	buffer, _ := io.ReadAll(res.Response().Body)
	body := string(buffer)

	//首先在小组页面中取得开始的cursor
	body = global.StrBetween(body, `nonAdminGroups":{"groups_tab":`, `},"adminGroups"`)

	startCursorJson, err2 := simplejson.NewJson([]byte(body))
	if err2 != nil {
		return false, nil
	}

	//取得第一条信息的cursor
	startCursor := startCursorJson.Get("tab_groups_list").Get("edges").GetIndex(0).Get("cursor").MustString()
	endCursor := startCursor

	//将第一条的cursor的值插入到groups中，因为后续循环直接就跳过第一个cursor了
	groups := make([]GroupsInfo, 0)
	name := startCursorJson.Get("tab_groups_list").Get("edges").GetIndex(0).Get("node").Get("name").MustString()
	id := startCursorJson.Get("tab_groups_list").Get("edges").GetIndex(0).Get("node").Get("id").MustString()
	url := startCursorJson.Get("tab_groups_list").Get("edges").GetIndex(0).Get("node").Get("url").MustString()
	cursor := startCursorJson.Get("tab_groups_list").Get("edges").GetIndex(0).Get("cursor").MustString()

	groups = append(groups, GroupsInfo{id, name, url, cursor})
	total := 1

	for len(endCursor) > 0 {
		variables := `{"count":10,"cursor":"` + endCursor + `","listType":"NON_ADMIN_MODERATOR_GROUPS","scale":1.5,"__relay_internal__pv__GroupsCometEntityMenuEmbeddedrelayprovider":false,"__relay_internal__pv__GroupsCometEntityMenuNotEmbeddedrelayprovider":true}`
		param.Fb_api_req_friendly_name = "GroupsLeftRailYourGroupsPaginatedQuery"
		param.Variables = variables
		param.Fb_api_caller_class = "RelayModern"
		param.Ccg = "EXCELLENT"
		param.Doc_id = doc_id

		success, result := graphql(param)
		if !success {
			groups = nil
			break
		}

		json, err := simplejson.NewJson([]byte(result))
		if err != nil {
			groups = nil
			break
		}

		node := json.Get("data").Get("viewer").Get("groups_tab").Get("tab_groups_list")
		endCursor = node.Get("page_info").Get("end_cursor").MustString()
		has_next_page := node.Get("page_info").Get("has_next_page").MustBool()
		length := len(node.Get("edges").MustArray())

		for i := 0; i < length; i++ {

			name := node.Get("edges").GetIndex(i).Get("node").Get("name").MustString()
			id := node.Get("edges").GetIndex(i).Get("node").Get("id").MustString()
			url := node.Get("edges").GetIndex(i).Get("node").Get("url").MustString()
			cursor := node.Get("edges").GetIndex(i).Get("cursor").MustString()

			groups = append(groups, GroupsInfo{id, name, url, cursor})

			total++
			if total >= num {
				endCursor = ""
				break
			}
		}

		//小组全部加载完毕
		if !has_next_page {
			endCursor = ""
		}
	}

	return len(groups) > 0, groups
}

func PostGroup(param GraphqlParam, doc_id string, groupId string, text string) (bool, string) {

	//首先检查这个小组是否允许用户发帖
	success, msg := checkGroupCanPost(param, "5417203551662631", groupId)
	if !success {
		return false, msg
	}

	milliTimestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	guid := uuid.New().String()
	variables := `{"input":{"composer_entry_point":"inline_composer","composer_source_surface":"group","composer_type":"group","logging":{"composer_session_id":"` + guid + `"},"source":"WWW","attachments":[],"message":{"ranges":[],"text":"` + text + `"},"with_tags_ids":[],"inline_activities":[],"explicit_place_id":"0","text_format_preset_id":"0","navigation_data":{"attribution_id_v2":"CometGroupDiscussionRoot.react,comet.group,unexpected,` + milliTimestamp + `,441376,2361831622,;GroupsCometCrossGroupFeedRoot.react,comet.groups.feed,tap_tabbar,` + milliTimestamp + `,822701,2361831622,"},"tracking":[null],"audience":{"to_id":"` + groupId + `"},"actor_id":"` + param.UserId + `","client_mutation_id":"` + global.RandomStr(2) + `"},"displayCommentsFeedbackContext":null,"displayCommentsContextEnableComment":null,"displayCommentsContextIsAdPreview":null,"displayCommentsContextIsAggregatedShare":null,"displayCommentsContextIsStorySet":null,"feedLocation":"GROUP","feedbackSource":0,"focusCommentID":null,"gridMediaWidth":null,"groupID":null,"scale":1.5,"privacySelectorRenderLocation":"COMET_STREAM","renderLocation":"group","useDefaultActor":false,"inviteShortLinkKey":null,"isFeed":false,"isFundraiser":false,"isFunFactPost":false,"isGroup":true,"isEvent":false,"isTimeline":false,"isSocialLearning":false,"isPageNewsFeed":false,"isProfileReviews":false,"isWorkSharedDraft":false,"UFI2CommentsProvider_commentsKey":"CometGroupDiscussionRootSuccessQuery","hashtag":null,"canUserManageOffers":false,"__relay_internal__pv__FBReelsEnableDeferrelayprovider":false}`

	param.Fb_api_req_friendly_name = "ComposerStoryCreateMutation"
	param.Variables = variables
	param.Fb_api_caller_class = "RelayModern"
	param.Ccg = "EXCELLENT"
	param.Doc_id = doc_id

	success, result := graphql(param)
	if !success {
		//为让社群免受垃圾信息打扰，你暂时无法发帖
		// if strings.Contains(result, `\u4e3a\u8ba9\u793e\u7fa4\u514d\u53d7\u5783\u573e\u4fe1\u606f\u6253\u6270`) {
		// 	return false, "为让社群免受垃圾信息打扰，你暂时无法发帖"
		// }

		//取出错误信息
		json, err := simplejson.NewJson([]byte(result))
		if err != nil {
			return false, "Json parse error" + "【" + result + "】"
		}

		summary := json.Get("errors").GetIndex(0).Get("summary").MustString()
		description_raw := json.Get("errors").GetIndex(0).Get("description_raw").MustString()

		return false, summary + description_raw + "!"
	}

	//发布成功，取出发布成功的帖子链接返回
	json, err := simplejson.NewJson([]byte(result))
	if err != nil {
		return false, "Json parse error" + "【" + result + "】"
	}

	url := json.Get("data").Get("story_create").Get("story").Get("url").MustString()
	if len(url) > 0 {
		return true, url
	}

	return false, result
}

func FriendRequest(param GraphqlParam, doc_id string, friendId string) (bool, string) {

	variables := `{"input":{"attribution_id_v2":"SearchCometGlobalSearchDefaultTabRoot.react,comet.search_results.default_tab,unexpected,1664676166332,428660,,;SearchCometGlobalSearchDefaultTabRoot.react,comet.search_results.default_tab,unexpected,1664676115989,939071,,;SearchCometGlobalSearchDefaultTabRoot.react,comet.search_results.default_tab,unexpected,1664676082266,173254,,;SearchCometGlobalSearchTopTabRoot.react,comet.search_results.top_tab,tap_search_bar,1664675178022,418344,391724414624676,","friend_requestee_ids":["` + friendId + `"],"refs":[null],"source":"search","warn_ack_for_ids":[],"actor_id":"` + param.UserId + `","client_mutation_id":"` + global.RandomStr(1) + `"},"scale":1.5}`

	param.Fb_api_req_friendly_name = "FriendingCometFriendRequestSendMutation"
	param.Variables = variables
	param.Fb_api_caller_class = "RelayModern"
	param.Ccg = "EXCELLENT"
	param.Doc_id = doc_id

	return graphql(param)
}

func GetAccountName(cookie string) (bool, string) {

	header := req.Header{
		`accept`:                    ` text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		`accept-language`:           ` zh-CN,zh;q=0.9`,
		`cache-control`:             ` no-cache`,
		`pragma`:                    ` no-cache`,
		`sec-ch-ua`:                 ` "Chromium";v="104", " Not A;Brand";v="99", "Google Chrome";v="104"`,
		`sec-ch-ua-mobile`:          ` ?0`,
		`sec-ch-ua-platform`:        ` "Windows"`,
		`sec-fetch-dest`:            ` document`,
		`sec-fetch-mode`:            ` navigate`,
		`sec-fetch-site`:            ` none`,
		`sec-fetch-user`:            ` ?1`,
		`upgrade-insecure-requests`: ` 1`,
		`cookie`:                    cookie,
		`user-agent`:                ` Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`,
	}

	res, err1 := req.Get("https://www.facebook.com/", header)
	if err1 != nil {
		return false, ""
	}

	buffer, _ := io.ReadAll(res.Response().Body)
	body := string(buffer)
	userId := global.StrBetween(cookie, `c_user=`, `;`)
	name := global.StrBetween(body, `["CurrentUserInitialData",[],{"ACCOUNT_ID":"`+userId+`","USER_ID":"`+userId+`","NAME":"`, `","SHORT_NAME"`)

	//将unicode字符转为可见字符
	str, _ := strconv.Unquote(`"` + name + `"`)

	if len(string(str)) > 0 {
		return true, str
	}

	return false, ""
}
