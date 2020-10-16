import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from './dash/dash.component';
import {MainComponent} from './main.component';
import {LinkComponent} from './link/link.component';
import {TunnelComponent} from './tunnel/tunnel.component';
import {PluginComponent} from './plugin/plugin.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {ModelComponent} from './model/model.component';
import {ModelAdapterComponent} from './model-detail/adapter/model-adapter.component';
import {ModelVariableComponent} from './model-detail/variable/model-variable.component';
import {ModelBatchComponent} from './model-detail/batch/model-batch.component';
import {ModelJobComponent} from './model-detail/job/model-job.component';
import {ModelStrategyComponent} from './model-detail/strategy/model-strategy.component';
import {ModelEditComponent} from './model-detail/model-edit/model-edit.component';
import {TunnelEditComponent} from './tunnel-edit/tunnel-edit.component';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {ModelStrategyEditComponent} from './model-detail/strategy-edit/model-strategy-edit.component';
import {ModelJobEditComponent} from './model-detail/job-edit/model-job-edit.component';
import {ModelBatchEditComponent} from './model-detail/batch-edit/model-batch-edit.component';
import {ModelVariableEditComponent} from './model-detail/variable-edit/model-variable-edit.component';
import {ModelAdapterEditComponent} from './model-detail/adapter-edit/model-adapter-edit.component';
import {PluginEditComponent} from './plugin-edit/plugin-edit.component';
import {ModelDetailComponent} from "./model-detail/model-detail.component";

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
      {path: 'model-detail/:id', component: ModelDetailComponent},
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
