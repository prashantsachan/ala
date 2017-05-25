package topology
import(
	"fmt"
	"errors"
	"github.com/samuel/go-zookeeper/zk"
    log "github.com/Sirupsen/logrus"
    "encoding/json"
)
const RootNode = "/topology"
var flags = int32(0)
var acl = zk.WorldACL(zk.PermAll)
type ZkServiceDao struct{
	Conn *zk.Conn
}

func (this *ZkServiceDao) Init(){
	// create if the root path doesn't exist
	exists,_,err:= this.Conn.Exists(RootNode)
	if err !=nil{
		log.WithFields(log.Fields{"module":"zkServiceDao","action":"exists",
			"path":RootNode,"error":err}).Error("unable to determine existance of root, assuming it exists")
	}else if !exists{
		path,err := this.Conn.Create(RootNode, []byte("[metricCollection]root node for topology"),flags, acl)
		if err == nil{
			log.WithFields(log.Fields{"module":"zkServiceDao","action":"create",
				"path":RootNode}).Info("created RootNode with path: "+path)
		}else{
			log.WithFields(log.Fields{"module":"zkServiceDao","action":"create",
				"path":RootNode, "error":err}).Info("unable to create RootNode")
		}

	}
}
// Fetches to get all The services, returns the first 
func (this *ZkServiceDao) GetAllServices()([]Service,error){
	// get all children of the RootNode
	ids,_,err := this.Conn.Children(RootNode)
	if err!=nil{
		return nil,err
	}else{
		var services []Service
		var failedIds []string
		for _,id := range ids{
			path := RootNode+"/"+id
			data,_,zkErr :=this.Conn.Get(path)
			if zkErr !=nil{
				log.WithFields(log.Fields{"module":"zkServiceDao","action":"getAll",
					"id":id,"error":zkErr}).Error("error fetching service from zk")
				failedIds = append(failedIds, id)
			}else{
				var s Service;
				jErr:= json.Unmarshal(data, &s)
				if jErr !=nil{
					log.WithFields(log.Fields{"module":"zkServiceDao","action":"getAll",
						"data":string(data),"error":zkErr}).Error("error parsing to Service")
					failedIds = append(failedIds, id)	
				}else{
					services = append(services, s)
				}
			}
		}
		if len(failedIds)==0 {
			return services, nil
		}else{
			msg:= fmt.Sprintf("failed to retrieve services for Ids: %v",failedIds)
			return services, errors.New(msg)
		}
	}
}

func (this *ZkServiceDao) AddService(s Service)error{
	// store this service as a child of root
	path := RootNode+"/"+s.Id
	data,jErr := json.Marshal(s)
	if jErr !=nil{
		return jErr
	}else{
		_, zkErr := this.Conn.Create(path, data, flags, acl)	
		return zkErr
	}
}
func (this *ZkServiceDao) DeleteService(id string) error{
	// delete this service
	path := RootNode+"/"+id
	zkErr := this.Conn.Delete(path, -1)
	return zkErr
	
}
