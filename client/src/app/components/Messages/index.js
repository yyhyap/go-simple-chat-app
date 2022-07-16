import React from 'react';

export class Messages extends React.Component {
    render() {
        const { messages } = this.props;

        return(            
            messages ? (
                messages.map(message =>
                    <p key={message.id}>
                        <strong>
                            {message.sender}:&nbsp;
                        </strong>
                        {message.body}
                    </p>
                )
            ) : (
                // same as <div></div>
                <div />
            )            
        )
    }
}