import { store } from "./store";

class MainView extends React.Component {
  render() {
    return (
    	<div>
    		<CreateBranch />
    		<hr/>
          	<FilterableBranchesTable branches={window.branches} />
    	</div>
    	)
  }
}

