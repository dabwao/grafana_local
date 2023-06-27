import { DataQuery as SchemaDataQuery, DataSourceRef as SchemaDataSourceRef } from '@grafana/schema';

import { AnalyzeQueryOptions, QueryFixAction } from './datasource';

/**
 * @deprecated use the type from @grafana/schema
 */
export interface DataQuery extends SchemaDataQuery {}

/**
 * @deprecated use the type from @grafana/schema
 */
export interface DataSourceRef extends SchemaDataSourceRef {}

/**
 * Attached to query results (not persisted)
 *
 * @public
 */
export enum DataTopic {
  Annotations = 'annotations',
}

/**
 * Abstract representation of any label-based query
 * @internal
 */
export interface AbstractQuery extends SchemaDataQuery {
  labelMatchers: AbstractLabelMatcher[];
}

/**
 * @internal
 */
export enum AbstractLabelOperator {
  Equal = 'Equal',
  NotEqual = 'NotEqual',
  EqualRegEx = 'EqualRegEx',
  NotEqualRegEx = 'NotEqualRegEx',
}

/**
 * @internal
 */
export type AbstractLabelMatcher = {
  name: string;
  value: string;
  operator: AbstractLabelOperator;
};

/**
 * @internal
 */
export interface DataSourceWithQueryImportSupport<TQuery extends SchemaDataQuery> {
  importFromAbstractQueries(labelBasedQuery: AbstractQuery[]): Promise<TQuery[]>;
}

/**
 * @internal
 */
export interface DataSourceWithQueryExportSupport<TQuery extends SchemaDataQuery> {
  exportToAbstractQueries(query: TQuery[]): Promise<AbstractQuery[]>;
}

/**
 * @internal
 */
export interface DataSourceWithQueryManipulationSupport<TQuery extends SchemaDataQuery> {
  /**
   * Used in explore
   */
  modifyQuery(query: TQuery, action: QueryFixAction): TQuery;

  /**
   * Used in explore for Log details
   *
   * @alpha
   */
  analyzeQuery?(query: TQuery, options: AnalyzeQueryOptions): boolean;
}

/**
 * @internal
 */
export const hasQueryImportSupport = <TQuery extends SchemaDataQuery>(
  datasource: unknown
): datasource is DataSourceWithQueryImportSupport<TQuery> => {
  return (datasource as DataSourceWithQueryImportSupport<TQuery>).importFromAbstractQueries !== undefined;
};

/**
 * @internal
 */
export const hasQueryExportSupport = <TQuery extends SchemaDataQuery>(
  datasource: unknown
): datasource is DataSourceWithQueryExportSupport<TQuery> => {
  return (datasource as DataSourceWithQueryExportSupport<TQuery>).exportToAbstractQueries !== undefined;
};

/**
 * @internal
 */
export const hasQueryManipulationSupport = <TQuery extends SchemaDataQuery>(
  datasource: unknown,
  method: keyof DataSourceWithQueryManipulationSupport<TQuery>
): datasource is DataSourceWithQueryManipulationSupport<TQuery> => {
  return Object.hasOwnProperty.call(datasource, method);
};
