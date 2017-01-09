import React, { Component } from 'react';
import './InputView.css'

class InputView extends Component {
  state = {
    input: '',
  }

  constructor(props) {
    super(props);

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleInput = this.handleInput.bind(this);
  }

  handleSubmit(event) {
    event.stopPropagation();
    event.preventDefault();

    this.props.onSubmit(this.state.input);
  }

  handleInput(event) {
    this.setState({input: event.target.value});
  }

  render() {
    return (
      <div className='flex-center'>
        <div className='input-view'>
          <h1>Dos</h1>

          <form onSubmit={this.handleSubmit}>
            <input
              type={this.props.type}
              placeholder={this.props.placeholder}
              onChange={this.handleInput}
              disabled={this.props.disabled}
              value={this.state.input} />

            <button
              disabled={this.state.input.length === 0 || this.props.disabled}>
                Join Game
            </button>
          </form>

          { this.props.error && <div className='error'>
            { this.props.error }
          </div> }

          { this.props.children }
        </div>
      </div>
    );
  }
}

export default InputView;
