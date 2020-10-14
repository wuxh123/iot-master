import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from './dash/dash.component';
import {MainComponent} from './main.component';
import {LinkComponent} from './link/link.component';
import {TunnelComponent} from './tunnel/tunnel.component';
import {PluginComponent} from './plugin/plugin.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {ModelComponent} from './model/model/model.component';
import {ModelTunnelComponent} from './model/tunnel/model-tunnel.component';
import {ModelVariableComponent} from './model/variable/model-variable.component';
import {ModelBatchComponent} from './model/batch/model-batch.component';
import {ModelJobComponent} from './model/job/model-job.component';
import {ModelStrategyComponent} from './model/strategy/model-strategy.component';
import {ModelEditComponent} from './model/model-edit/model-edit.component';
import {TunnelEditComponent} from './tunnel-edit/tunnel-edit.component';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {ModelStrategyEditComponent} from './model/strategy-edit/model-strategy-edit.component';
import {ModelJobEditComponent} from './model/job-edit/model-job-edit.component';
import {ModelBatchEditComponent} from './model/batch-edit/model-batch-edit.component';
import {ModelVariableEditComponent} from './model/variable-edit/model-variable-edit.component';
import {ModelTunnelEditComponent} from './model/tunnel-edit/model-tunnel-edit.component';
import {PluginEditComponent} from './plugin-edit/plugin-edit.component';

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
      {path: 'model', component: ModelComponent},
      {path: 'model-create', component: ModelEditComponent},
      {path: 'model-edit/:id', component: ModelEditComponent},
      {path: 'model-tunnel', component: ModelTunnelComponent},
      {path: 'model-tunnel-create', component: ModelTunnelEditComponent},
      {path: 'model-tunnel-edit/:id', component: ModelTunnelEditComponent},
      {path: 'model-variable', component: ModelVariableComponent},
      {path: 'model-variable-create', component: ModelVariableEditComponent},
      {path: 'model-variable-edit/:id', component: ModelVariableEditComponent},
      {path: 'model-batch', component: ModelBatchComponent},
      {path: 'model-batch-create', component: ModelBatchEditComponent},
      {path: 'model-batch-edit/:id', component: ModelBatchEditComponent},
      {path: 'model-job', component: ModelJobComponent},
      {path: 'model-job-create', component: ModelJobEditComponent},
      {path: 'model-job-edit/:id', component: ModelJobEditComponent},
      {path: 'model-strategy', component: ModelStrategyComponent},
      {path: 'model-strategy-create', component: ModelStrategyEditComponent},
      {path: 'model-strategy-edit/:id', component: ModelStrategyEditComponent},
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
