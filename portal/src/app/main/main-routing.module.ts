import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from './dash/dash.component';
import {MainComponent} from './main.component';
import {LinkComponent} from './link/link.component';
import {TunnelComponent} from './tunnel/tunnel.component';
import {PluginComponent} from './plugin/plugin.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {ProjectComponent} from './project/project.component';
import {ModelAdapterComponent} from './project-adapter/model-adapter.component';
import {ModelVariableComponent} from './project-variable/model-variable.component';
import {ModelBatchComponent} from './project-batch/model-batch.component';
import {ModelJobComponent} from './project-job/model-job.component';
import {ModelStrategyComponent} from './project-strategy/model-strategy.component';
import {ProjectEditComponent} from './project-edit/project-edit.component';
import {TunnelEditComponent} from './tunnel-edit/tunnel-edit.component';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {ModelStrategyEditComponent} from './project-strategy-edit/model-strategy-edit.component';
import {ModelJobEditComponent} from './project-job-edit/model-job-edit.component';
import {ModelBatchEditComponent} from './project-batch-edit/model-batch-edit.component';
import {ModelVariableEditComponent} from './project-variable-edit/model-variable-edit.component';
import {ModelAdapterEditComponent} from './project-adapter-edit/model-adapter-edit.component';
import {PluginEditComponent} from './plugin-edit/plugin-edit.component';
import {ProjectDetailComponent} from './project-detail/project-detail.component';

const routes: Routes = [
  {
    path: '',
    component: MainComponent,
    children: [
      {path: '', redirectTo: 'dash'},
      {path: 'dash', component: DashComponent},
      {path: 'tunnel', component: TunnelComponent},
      {path: 'tunnel/create', component: TunnelEditComponent},
      {path: 'tunnel/:id/edit', component: TunnelEditComponent},
      {path: 'tunnel/:id/link', component: LinkComponent},
      {path: 'link', component: LinkComponent},
      {path: 'link/:id/edit', component: LinkEditComponent},
      {path: 'link/:id/monitor', component: LinkMonitorComponent},
      {path: 'plugin', component: PluginEditComponent},
      {path: 'plugin/create', component: PluginEditComponent},
      {path: 'plugin/:id/edit', component: PluginComponent},
      {path: 'project', component: ProjectComponent},
      {path: 'project/create', component: ProjectEditComponent},
      {path: 'project/:id/edit', component: ProjectEditComponent},
      {path: 'project/:id/detail', component: ProjectDetailComponent},
      {path: 'project/:id/adapter', component: ModelAdapterComponent},
      {path: 'project/:id/adapter/create', component: ModelAdapterEditComponent},
      {path: 'project/:id/adapter/:id/edit', component: ModelAdapterEditComponent},
      {path: 'project/:id/variable', component: ModelVariableComponent},
      {path: 'project/:id/variable/create', component: ModelVariableEditComponent},
      {path: 'project/:id/variable/:id/edit', component: ModelVariableEditComponent},
      {path: 'project/:id/batch', component: ModelBatchComponent},
      {path: 'project/:id/batch/create', component: ModelBatchEditComponent},
      {path: 'project/:id/batch/:id/edit', component: ModelBatchEditComponent},
      {path: 'project/:id/job', component: ModelJobComponent},
      {path: 'project/:id/job/create', component: ModelJobEditComponent},
      {path: 'project/:id/job/:id/edit', component: ModelJobEditComponent},
      {path: 'project/:id/strategy', component: ModelStrategyComponent},
      {path: 'project/:id/strategy/create', component: ModelStrategyEditComponent},
      {path: 'project/:id/strategy/:id/edit', component: ModelStrategyEditComponent},

      {path: '**', redirectTo: 'dash'},
    ]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MainRoutingModule {
}
