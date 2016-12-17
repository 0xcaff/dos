import React, { Component } from 'react';
import './SocketStatus.css'

class SocketStatus extends Component {
  state = {};

  constructor(props) {
    super(props);

    this.handleChange = this.handleChange.bind(this);

    // eslint-disable-next-line react/no-direct-mutation-state
    this.state.readyState = props.socket.readyState;

    props.socket.addEventListener('close', this.handleChange);
    props.socket.addEventListener('open', this.handleChange);
  }

  handleChange() {
    this.setState({
      readyState: this.props.socket.readyState,
    });
  }

  render() {
    if (this.state.readyState === 0) {
      // Connecting
      return (<div className='connecting status'>
        { 'Connecting ' }
        <span className='loading'>
          <span>.</span>
          <span>.</span>
          <span>.</span>
        </span>
      </div>)
    } else if (this.state.readyState === 2 || this.state.readyState === 3) {
      // Closing or Closed
      return <div className='disconnected status'>Disconnected</div>
    } else {
      return null;
    }
  }
}

export default SocketStatus;
