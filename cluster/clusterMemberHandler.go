package cluster

// "encoding/json"
import (
	"net/http"

	gmux "github.com/gorilla/mux"
	util "github.com/tmax-cloud/hypercloud-api-server/util"
	caller "github.com/tmax-cloud/hypercloud-api-server/util/caller"
	clusterDataFactory "github.com/tmax-cloud/hypercloud-api-server/util/dataFactory/cluster"
	"k8s.io/klog"
	// "encoding/json"
)

func ListClusterMember(res http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	userId := queryParams.Get(QUERY_PARAMETER_USER_ID)
	userGroups := queryParams[util.QUERY_PARAMETER_USER_GROUP]
	vars := gmux.Vars(req)
	cluster := vars["clustermanager"]
	namespace := vars["namespace"]

	if err := util.StringParameterException(userGroups, userId, cluster, namespace); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusBadRequest)
		return
	}

	// cluster ready 인지 확인
	// var clm *clusterv1alpha1.ClusterManager
	clm, err := caller.GetCluster(userId, userGroups, cluster, namespace)
	if err != nil {
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if !clm.Status.Ready || clm.Status.Phase == "Deleting" {
		msg := "Cannot invite member to cluster: cluster is deleting or not ready"
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, nil, http.StatusBadRequest)
		return
	}

	if userId == clm.Annotations["owner"] {
		clusterMemberList, err := clusterDataFactory.ListClusterMember(cluster, namespace)
		if err != nil {
			klog.V(1).Infoln(err)
			util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
			return
		}
		msg := "List cluster success"
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, clusterMemberList, http.StatusOK)
		return
	} else {
		clusterMemberList, err := clusterDataFactory.ListClusterMemberWithOutPending(cluster, namespace)
		if err != nil {
			klog.V(1).Infoln(err)
			util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
			return
		}
		msg := "List cluster success"
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, clusterMemberList, http.StatusOK)
		return
	}
}

func ListClusterMemberWithOutPending(res http.ResponseWriter, req *http.Request) {
	vars := gmux.Vars(req)
	cluster := vars["clustermanager"]
	namespace := vars["namespace"]
	klog.V(3).Infoln("in func ListClusterMemberWithOutPending")
	klog.V(3).Infoln("cluster = " + cluster + ", namespace = " + namespace)

	clusterMemberList, err := clusterDataFactory.ListClusterMemberWithOutPending(cluster, namespace)
	if err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	msg := "List cluster success"
	klog.V(3).Infoln(msg)
	util.SetResponse(res, msg, clusterMemberList, http.StatusOK)
}

func ListClusterInvitedMember(res http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	userId := queryParams.Get(QUERY_PARAMETER_USER_ID)
	userGroups := queryParams[util.QUERY_PARAMETER_USER_GROUP]
	vars := gmux.Vars(req)
	cluster := vars["clustermanager"]
	namespace := vars["namespace"]

	if err := util.StringParameterException(userGroups, userId, cluster, namespace); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusBadRequest)
		return
	}

	clm, err := caller.GetCluster(userId, userGroups, cluster, namespace)
	if err != nil {
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if !clm.Status.Ready || clm.Status.Phase == "Deleting" {
		msg := "Cannot list invited member in cluster: cluster is deleting or not ready"
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, nil, http.StatusBadRequest)
		return
	}

	clusterMemberList, err := clusterDataFactory.ListClusterInvitedMember(cluster, namespace)
	if err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	msg := "List cluster invited member success"
	klog.V(3).Infoln(msg)
	util.SetResponse(res, msg, clusterMemberList, http.StatusOK)
}

func ListClusterGroup(res http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	userId := queryParams.Get(QUERY_PARAMETER_USER_ID)
	userGroups := queryParams[util.QUERY_PARAMETER_USER_GROUP]
	vars := gmux.Vars(req)
	cluster := vars["clustermanager"]
	namespace := vars["namespace"]

	clm, err := caller.GetCluster(userId, userGroups, cluster, namespace)
	if err != nil {
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if !clm.Status.Ready || clm.Status.Phase == "Deleting" {
		msg := "Cannot list invited member in cluster: cluster is deleting or not ready"
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, nil, http.StatusBadRequest)
		return
	}

	clusterMemberList, err := clusterDataFactory.ListClusterGroup(cluster, namespace)
	if err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	msg := "List cluster group success"
	klog.V(3).Infoln(msg)
	util.SetResponse(res, msg, clusterMemberList, http.StatusOK)
}

// cluster에서 초대 받은 member(User/Group) 제거시 호출
func RemoveMember(res http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	userId := queryParams.Get(QUERY_PARAMETER_USER_ID)
	userGroups := queryParams[util.QUERY_PARAMETER_USER_GROUP]

	vars := gmux.Vars(req)
	cluster := vars["clustermanager"]
	attribute := vars["attribute"]
	memberId := vars["member"]
	namespace := vars["namespace"]

	if err := util.StringParameterException(userGroups, userId, cluster, attribute, memberId, namespace); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusBadRequest)
		return
	}

	clm, err := caller.GetCluster(userId, userGroups, cluster, namespace)
	if err != nil {
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	if !clm.Status.Ready || clm.Status.Phase == "Deleting" {
		msg := "Cannot remove member in cluster: cluster is deleting phase or not ready"
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, nil, http.StatusBadRequest)
		return
	}

	clusterMember := util.ClusterMemberInfo{}
	clusterMember.Namespace = namespace
	clusterMember.Cluster = cluster
	clusterMember.MemberId = memberId
	clusterMember.Attribute = attribute
	clusterMember.Status = "invited"

	clusterMemberList, err := clusterDataFactory.ListClusterMember(clusterMember.Cluster, clusterMember.Namespace)
	if err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	var clusterOwner string
	var existMember []string
	for _, val := range clusterMemberList {
		if val.Status == "owner" {
			clusterOwner = val.MemberId
		} else {
			existMember = append(existMember, val.MemberId)
		}
	}

	if userId != clusterOwner {
		msg := "Request user [ " + userId + " ]is not a cluster owner [ " + clusterOwner + " ]"
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, nil, http.StatusBadRequest)
		return
	}

	if !util.Contains(existMember, memberId) {
		msg := attribute + " [ " + memberId + " ] is already removed in cluster [ " + cluster + " ] "
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, nil, http.StatusBadRequest)
		return
	}

	// db에서 삭제
	if err := clusterDataFactory.Delete(clusterMember); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	// role 삭제
	if err := caller.RemoveRoleFromRemote(clm, memberId, attribute); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if err := caller.RemoveSASecretFromRemote(clm, memberId, attribute); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if err := caller.RemoveRemoteSecretInLocal(clm, memberId, attribute); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if err := caller.DeleteCLMRole(clm, memberId, attribute); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	if err := caller.DeleteNSGetRole(clm, memberId, attribute); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	msg := ""
	if attribute == "user" {
		msg = "User [" + memberId + "] is removed from cluster [" + clm.Name + "]"
	} else {
		msg = "Group [" + memberId + "] is removed from cluster [" + clm.Name + "]"
	}

	klog.V(3).Infoln(msg)
	util.SetResponse(res, msg, nil, http.StatusOK)
}

// cluster에서 초대 받은 member의 역할 변경시 호출
func UpdateMemberRole(res http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	userId := queryParams.Get(QUERY_PARAMETER_USER_ID)
	userGroups := queryParams[util.QUERY_PARAMETER_USER_GROUP]
	remoteRole := queryParams.Get(QUERY_PARAMETER_REMOTE_ROLE)

	vars := gmux.Vars(req)
	cluster := vars["clustermanager"]
	attribute := vars["attribute"]
	memberId := vars["member"]
	namespace := vars["namespace"]

	if err := util.StringParameterException(userGroups, userId, cluster, attribute, memberId, remoteRole, namespace); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusBadRequest)
		return
	}

	clm, err := caller.GetCluster(userId, userGroups, cluster, namespace)
	if err != nil {
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	if !clm.Status.Ready || clm.Status.Phase == "Deleting" {
		msg := "Cannot invite member to cluster in deleting phase or not ready status"
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, nil, http.StatusBadRequest)
		return
	}

	clusterMember := util.ClusterMemberInfo{}
	clusterMember.Namespace = namespace
	clusterMember.Cluster = cluster
	clusterMember.MemberId = memberId
	clusterMember.Role = remoteRole
	clusterMember.Attribute = attribute
	clusterMember.Status = "invited"

	clusterMemberList, err := clusterDataFactory.ListClusterMember(clusterMember.Cluster, clusterMember.Namespace)
	if err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	var clusterOwner string
	var existMember []string
	for _, val := range clusterMemberList {
		if val.Status == "owner" {
			clusterOwner = val.MemberId
		} else {
			existMember = append(existMember, val.MemberId)
		}
	}

	if userId != clusterOwner {
		msg := "Request user [ " + userId + " ]is not a cluster owner [ " + clusterOwner + " ]"
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, nil, http.StatusBadRequest)
		return
	}

	if !util.Contains(existMember, memberId) {
		msg := attribute + " [ " + memberId + " ] is not in cluster [ " + cluster + " ] "
		klog.V(3).Infoln(msg)
		util.SetResponse(res, msg, nil, http.StatusBadRequest)
		return
	}

	// db에서 role update
	if err := clusterDataFactory.UpdateRole(clusterMember); err != nil {
		klog.V(1).Infoln(err)
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	// role 삭제 후 재 생성
	if err := caller.RemoveRoleFromRemote(clm, memberId, attribute); err != nil {
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if err := caller.CreateRoleInRemote(clm, memberId, remoteRole, clusterMember.Attribute); err != nil {
		util.SetResponse(res, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	msg := attribute + " [" + memberId + "] role is updated to [" + remoteRole + "] in cluster [" + clm.Name + "]"
	klog.V(3).Infoln(msg)
	util.SetResponse(res, msg, nil, http.StatusOK)
}
