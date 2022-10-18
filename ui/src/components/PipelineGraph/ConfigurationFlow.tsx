import { useCallback, useEffect } from "react";
import ReactFlow, { Controls, useReactFlow } from "react-flow-renderer";
import SourceNode from "./Nodes/SourceNode";
import DestinationNode from "./Nodes/DestinationNode";
import UIControlNode from "./Nodes/UIControlNode";
import {
  getNodesAndEdges,
  GRAPH_NODE_OFFSET,
  GRAPH_PADDING,
  Page,
  updateMetricData,
} from "../../utils/graph/utils";
import { ProcessorNode } from "./Nodes/ProcessorNode";
import { useConfigurationPage } from "../../pages/configurations/configuration/ConfigurationPageContext";
import { gql } from "@apollo/client";
import { useConfigurationMetricsSubscription } from "../../graphql/generated";
import OverviewEdge from '../../pages/overview/OverviewEdge';
import ConfigurationEdge from './Nodes/ConfigurationEdge';

gql`
  subscription ConfigurationMetrics($period: String!, $name: String!) {
    configurationMetrics(period: $period, name: $name) {
      metrics {
        name
        nodeID
        pipelineType
        value
        unit
      }
    }
  }
`;

const nodeTypes = {
  sourceNode: SourceNode,
  destinationNode: DestinationNode,
  uiControlNode: UIControlNode,
  processorNode: ProcessorNode,
};

const edgeTypes = {
  overviewEdge: OverviewEdge,
  configurationEdge: ConfigurationEdge,
};

export const TARGET_OFFSET_MULTIPLIER = 250;

interface ConfigurationFlowProps {
  period: string;
  selectedTelemetry: string;
}

export const ConfigurationFlow: React.FC<ConfigurationFlowProps> = ({
  period,
  selectedTelemetry,
}) => {
  const reactFlowInstance = useReactFlow();
  const { configuration } = useConfigurationPage();
  const { nodes, edges } = getNodesAndEdges(
    Page.Configuration,
    configuration.graph!,
    TARGET_OFFSET_MULTIPLIER
  );

  const { data } = useConfigurationMetricsSubscription({
    variables: {
      period,
      name: configuration.metadata.name,
    },
  });

  updateMetricData(
    Page.Configuration,
    nodes,
    edges,
    data?.configurationMetrics.metrics ?? [],
    period,
    selectedTelemetry
  );

  const viewPortHeight =
    GRAPH_PADDING +
    Math.max(
      configuration.graph?.sources?.length || 0,
      configuration.graph?.targets?.length || 0
    ) *
      GRAPH_NODE_OFFSET;
  const onNodesChange = useCallback(() => {
    reactFlowInstance.fitView();
  }, [reactFlowInstance]);

  useEffect(() => {
    setTimeout(() => {
      reactFlowInstance.fitView();
    }, 10);
  }, [reactFlowInstance, viewPortHeight]);

  return (
    <div style={{ height: viewPortHeight, width: "100%" }}>
      <ReactFlow
        defaultNodes={nodes}
        defaultEdges={edges}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        proOptions={{account: 'paid-pro', hideAttribution: true}}
        nodesConnectable={false}
        nodesDraggable={false}
        fitView={true}
        deleteKeyCode={null}
        zoomOnScroll={false}
        panOnDrag={false}
        minZoom={0.1}
        onWheel={(event) => {
          window.scrollBy(event.deltaX, event.deltaY);
        }}
        // This callback only fires when a node is deleted in the graph (ie select node and press the delete button)
        // Need to figure out how to update after node is deleted from the App
        // This callback *does* fire when a Node is added in the App (?)
        onNodesChange={onNodesChange}
      >
        <Controls showZoom={false} showInteractive={false} />
      </ReactFlow>
    </div>
  );
};
