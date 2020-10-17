import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from './dash/dash.component';
import {MainComponent} from './main.component';
import {LinkComponent} from './link/link.component';
import {TunnelComponent} from './tunnel/tunnel.component';
import {PluginComponent} from './plugin/plugin.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {ProjectComponent} from './project/project.component';
import {ModelAdapterComponent} from './project-detail/adapter/model-adapter.component';
import {ModelVariableComponent} from './project-detail/variable/model-variable.component';
import {ModelBatchComponent} from './project-detail/batch/model-batch.component';
import {ModelJobComponent} from './project-detail/job/model-job.component';
import {ModelStrategyComponent} from './project-detail/strategy/model-strategy.component';
import {ProjectEditComponent} from './project-edit/project-edit.component';
import {TunnelEditComponent} from './tunnel-edit/tunnel-edit.component';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {ModelStrategyEditComponent} from './project-detail/strategy-edit/model-strategy-edit.component';
import {ModelJobEditComponent} from './project-detail/job-edit/model-job-edit.component';
import {ModelBatchEditComponent} from './project-detail/batch-edit/model-batch-edit.component';
import {ModelVariableEditComponent} from './project-detail/variable-edit/model-variable-edit.component';
import {ModelAdapterEditComponent} from './project-detail/adapter-edit/model-adapter-edit.component';
import {PluginEditComponent} from './plugin-edit/plugin-edit.component';
import {ProjectDetailComponent} from "./project-detail/project-detail.component";

const routes: Routes = [
  {
    path: '',
    component: MainComponent,
    children: [
      {path: '', redirectTo: 'dash'},
      {path: 'dash', component: DashComponent},
      {path: 'tunnel', component: TunnelComponent},
      {path: 'tunnel-create', component: TunnelEditComponent},
      {path: 'tunnel-edit/:id', component: TunnelEditComponent},
      {path: 'link', component: LinkComponent},
      {path: 'link-edit/:id', component: LinkEditComponent},
      {path: 'link-monitor/:id', component: LinkMonitorComponent},
      {path: 'plugin', component: PluginEditComponent},
      {path: 'plugin-create', component: PluginEditComponent},
      {path: 'plugin-edit/:id', component: PluginComponent},
      {path: 'model', component: ProjectComponent},
      {path: 'project-create', component: ProjectEditComponent},
      {path: 'project-edit/:id', component: ProjectEditComponent},
      {path: 'project-detail/:id', component: ProjectDetailComponent},
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
