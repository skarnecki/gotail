window.App = ((socket) => {
    class ContentRow extends React.Component {
        _clicked(e) {
            //TODO use props and selected
            const classess = e.target.className
            if (classess.indexOf("info") === -1) {
                e.target.className += "info"
            } else {
                e.target.className = classess.replace("info", "")
            }
        }

        render() {
            const { index, data, selected } = this.props
            const rowClass = selected ? "clickable-row info" : "clickable-row"
            return <tr key={`row-${index}`} onClick={ (e) => this._clicked(e) } className={ rowClass }>
                <td key={`row-cell-${index}`}>{ data }</td>
            </tr>;
        }
    }

    class FileContent extends React.Component {
      render() {
        const { lines, clicked, dispatch } = this.props
        const rows = lines.map((line, index) => {
            return <ContentRow index={ index } data={ line.data } selected= { line.selected }/>
        });
        return <table className="table table-hover"><tbody>
            {rows}
        </tbody></table>
      }
    }

    FileContent.propTypes = {
      lines: React.PropTypes.array.isRequired
    }

    function fileContents(state, action) {
        if (typeof state === 'undefined') {
          return {lines: []}
        }
        switch (action.type) {
          case 'APPEND':
            state.lines.push({data: action.lines, selected: false})
            return state
          case 'MARK_LINE':
            return state
          default:
            return state
        }
    }

    function render() {
        ReactDOM.render(
            <FileContent lines={ store.getState().lines } dispatch= { store.dispatch }/>,
            document.getElementById('content')
        )
    }

    const store = Redux.createStore(fileContents)
    render()
    store.subscribe(render)
    socket.onmessage = (msg) => {
        if (msg.data) {
            store.dispatch({ type: 'APPEND', lines: msg.data })
        }
    }
})
