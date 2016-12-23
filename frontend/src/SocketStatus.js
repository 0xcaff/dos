import React, { Component } from 'react';
import './SocketStatus.css'

class SocketStatus extends Component {
  render() {
    if (this.props.readyState === 0) {
      // Connecting
      return (<div className='connecting status'>
        { 'Connecting ' }
        <span className='loading'>
          <span>.</span>
          <span>.</span>
          <span>.</span>
        </span>
      </div>)
    } else if (this.props.readyState === 2 || this.props.readyState === 3) {
      // Closing or Closed
      return <div className='disconnected status'>Disconnected</div>
    } else {
      return null;
    }
  }
}

export default SocketStatus;
