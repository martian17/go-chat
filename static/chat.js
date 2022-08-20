class WS extends WebSocket{
    opened = false;
    pre_queue = [];
    constructor(){
        super(...arguments);
        let that = this;
        this.once("open",()=>{
            that.opened = true;
            that.pre_queue.map(res=>1);
        });
        
    }
    async send(msg){
        let that = this;
        if(!this.opened){
            await new Promise((res,rej)=>{
                that.pre_queue.push(res);
            });
        }
        super.send(msg);
        return msg;
    }
    once(ename,cb){
        let that = this;
        this.addEventListener(ename,()=>{
            this.removeEventListener(ename,cb);
            cb();
        });
    }
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




class MsgContainer extends ELEM{
    constructor(){
        super("div","class:msg-container");
        this.head = super.add("div","class:filler");
    }
    MAX_SIZE = 30;
    add(type,str){
        if(this.children.size >= this.MAX_SIZE){
            //remove one element from the front
            this.children.getNext(this.head).remove();
        }
        if(type === 0){
            super.add("div","class:msg self").add("div",0,str);
        }else if(type === 1){
            super.add("div","class:msg other").add("div",0,str);
        }
        //gotta scroll
        this.scrollBottom();
    }
    scrollBottom(){
        this.e.scroll(0,1000000);
    }
};

class MsgInput extends ELEM{
    constructor(chat){
        super("div","class:input");
        let input = this.add("input");
        input.on("keydown",(e)=>{
            let val = input.e.value.trim();
            if(e.key !== "Enter" || e.shiftKey || val === "")return;
            e.preventDefault();
            let msgval = val.replace(/\n/g,"<br>");
            chat.send(msgval);
            input.e.value = "";
        });
    }
};


class Chat extends ELEM{
    constructor(){
        super("div","class:chat");
        let that = this;
        //setup websocket first
        let ws = new WS("ws://"+location.host+"/socket");
        
        ws.on("message",(e)=>{
            console.log(e.data);
            that.receive(e.data);
        });
        this.ws = ws;
        
        
        let centerWindow = this.add("div");
        this.messages = new MsgContainer();
        let upper = centerWindow.add(this.messages);
        let lower = centerWindow.add(new MsgInput(this));
    }
    send(str){
        this.messages.add(0,str);
        this.ws.send(str);
    }
    receive(str){
        this.messages.add(1,str);
    }
};