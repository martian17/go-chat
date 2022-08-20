let body = new ELEM(document.body);

class WS extends WebSocket{
    on(ename,cb){
        let that = this;
        this.addEventListener(ename,cb);
        return {
            remove:()=>{
                that.removeEventListener(name,cb);
            }
        }
    }    
};



let ws = new WS("ws://"+location.host+"/socket");
ws.on("open",()=>{
    //p1.setInner("open");
});


ws.on("message",(e)=>{
    //p2.setInner(e.data);
});



//defining some API object
const API = {
    //returns promise
    get(path,content_type){
        return fetch(location.origin+path,{
            method: 'POST',
            headers: {
                'Content-Type': content_type
            },
            body:body
        });
    },
    
    //returns promise
    post(path,body,content_type){
        return fetch(location.origin+path,{
            method: 'POST',
            headers: {
                'Content-Type': content_type
            },
            body:body
        });
    },
    postBlob(path,body){
        return this.post(path,body,'application/octet-stream');
    },
    postJSON(path,body){
        return this.post(path,body,'application/json; charset=UTF-8');
    },
    postText(path,body){
        return this.post(path,body,'text/plain; charset=UTF-8');
    }
};


//test code
let testblob = new Blob([new Uint8Array([1,1,2,3,5,8])]);
/*
await API.postBlob("/api/post_buffer",testblob);
await API.postText("/api/post_string","from client");
let r = await API.get("/api/get_buffer");
console.log(new Uint8Array(await r.arrayBuffer()));
let r = await API.get("/api/get_string");
console.log(await r.text());
*/








