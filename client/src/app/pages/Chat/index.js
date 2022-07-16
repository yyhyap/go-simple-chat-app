import React from 'react';
import { InputText } from '../../components/InputText';
import { Messages } from '../../components/Messages';
import { Status } from '../../components/Status';
import './styles.css';

const baseURL = 'ws://localhost:8080/chat';

export class Chat extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            ws: undefined,
            username: '',
            message: '',
            messages: []
        }
    }
    
    render() {
        // ws from this.state
        const { ws, messages } = this.state;

        return (
            <div className='chat'>
                <h1>WebChat</h1>
                <Status status={ws ? 'connected' : 'disconnected'} />
                {
                    ws ? (<Messages messages={messages}/>) : (<div />)
                }
                <div className='chat-inputs'>
                    <InputText 
                        placeholder = {ws ? 'Enter your message' : 'Enter with your username'}
                        onChange = {value => ws ? this.setMessage(value) : this.setUsername(value)}
                        defaultValue = {ws ? this.state.message : this.state.username}
                    />
                </div>
                <button type='button' onClick={() => {ws ? this.sendMessage() : this.enterChat()}}>
                    {ws ? 'Send' : 'Enter'}
                </button>
            </div>
        )
    }

    enterChat() {
        const { username } = this.state;

        let ws = new WebSocket(baseURL + `?username=${username}`);

        ws.onopen = (e) => {
            console.log('Websocket connection established!', {e});
        }

        ws.onclose = (e) => {
            console.log('Websocket connection closed!', {e});
        }

        ws.onmessage = (msg) => {
            console.log('Websocket message: ', {msg});
            this.setMessages(msg.data)
        }

        ws.onerror = (error) => {
            console.log('Websocket error: ', {error});
            this.setState({ ws: undefined })
        }

        this.setState({ ws, username: '' }, () => {
            console.log('state after entered chat: ',this.state);
        });
        
    }

    sendMessage() {
        const { ws, message } = this.state;

        ws.send(message);
        this.setMessage('');
    }

    setUsername(value) {
        this.setState({username: value});
    }

    setMessage(value) {
        this.setState({message: value});
    }

    setMessages(value) {
        let messages = this.state.messages.concat([JSON.parse(value)]);
        this.setState({ messages });
    }
}

