import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from './dash/dash.component';
import {MainComponent} from './main.component';
import {LinkComponent} from './link/link.component';
import {ChannelComponent} from './channel/channel.component';
import {PluginComponent} from './plugin/plugin.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {ModelComponent} from './model/model.component';
import {TunnelComponent} from './tunnel/tunnel.component';
import {VariableComponent} from './variable/variable.component';
import {BatchComponent} from './batch/batch.component';
import {JobComponent} from './job/job.component';
import {StrategyComponent} from './strategy/strategy.component';
import {ModelEditComponent} from './model-edit/model-edit.component';
import {ChannelEditComponent} from './channel-edit/channel-edit.component';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {StrategyEditComponent} from './strategy-edit/strategy-edit.component';
import {JobEditComponent} from './job-edit/job-edit.component';
import {BatchEditComponent} from './batch-edit/batch-edit.component';
import {VariableEditComponent} from './variable-edit/variable-edit.component';
import {TunnelEditComponent} from './tunnel-edit/tunnel-edit.component';
import {PluginEditComponent} from './plugin-edit/plugin-edit.component';

const routes: Routes = [
  {
    path: '',
    component: MainComponent,
    children: [
      {path: '', redirectTo: 'dash'},
      {path: 'dash', component: DashComponent},
      {path: 'channel', component: ChannelComponent},
      {path: 'channel-create', component: ChannelEditComponent},
      {path: 'channel-edit/:id', component: ChannelEditComponent},
      {path: 'link', component: LinkComponent},
      {path: 'link-edit/:id', component: LinkEditComponent},
      {path: 'link-monitor/:id', component: LinkMonitorComponent},
      {path: 'plugin', component: PluginEditComponent},
      {path: 'plugin-create', component: PluginEditComponent},
      {path: 'plugin-edit/:id', component: PluginComponent},
      {path: 'model', component: ModelComponent},
      {path: 'model-create', component: ModelEditComponent},
      {path: 'model-edit/:id', component: ModelEditComponent},
      {path: 'tunnel', component: TunnelComponent},
      {path: 'tunnel-create', component: TunnelEditComponent},
      {path: 'tunnel-edit/:id', component: TunnelEditComponent},
      {path: 'variable', component: VariableComponent},
      {path: 'variable-create', component: VariableEditComponent},
      {path: 'variable-edit/:id', component: VariableEditComponent},
      {path: 'batch', component: BatchComponent},
      {path: 'batch-create', component: BatchEditComponent},
      {path: 'batch-edit/:id', component: BatchEditComponent},
      {path: 'job', component: JobComponent},
      {path: 'job-create', component: JobEditComponent},
      {path: 'job-edit/:id', component: JobEditComponent},
      {path: 'strategy', component: StrategyComponent},
      {path: 'strategy-create', component: StrategyEditComponent},
      {path: 'strategy-edit/:id', component: StrategyEditComponent},
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
