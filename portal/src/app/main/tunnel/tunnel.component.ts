import {Component, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzModalService, NzTableQueryParams} from 'ng-zorro-antd';
import {Router} from '@angular/router';
import {TabRef} from '../tabs/tabs.component';
import {TunnelEditComponent} from '../tunnel-edit/tunnel-edit.component';

@Component({
  selector: 'app-tunnel',
  templateUrl: './tunnel.component.html',
  styleUrls: ['./tunnel.component.scss']
})
export class TunnelComponent implements OnInit {

  channels: any[];
  total = 0;
  pageIndex = 1;
  pageSize = 10;
  sortField = null;
  sortOrder = null;
  filters = [];
  keyword = '';
  loading = false;

  netFilters = [
    {text: '串口', value: 'serial'},
    {text: 'TCP服务端', value: 'tcp-server'},
    {text: 'TCP客户端', value: 'tcp-client'},
    {text: 'UDP服务端', value: 'udp-server'},
    {text: 'UDP客户端', value: 'udp-client'},
  ];
  statusFilters = [{text: '启动', value: true}];


  constructor(private as: ApiService, private router: Router, private tab: TabRef, private ms: NzModalService) {
    tab.name = '通道管理';
  }

  ngOnInit(): void {
    this.loadFilters();
  }

  reload(): void {
    this.pageIndex = 1;
    this.keyword = '';
    this.load();
  }

  load(): void {
    this.loading = true;
    this.as.post('tunnels', {
      offset: (this.pageIndex - 1) * this.pageSize,
      length: this.pageSize,
      sortKey: this.sortField,
      sortOrder: this.sortOrder,
      filters: this.filters,
      keywords: [
        {key: 'Name', value: this.keyword},
        {key: 'Type', value: this.keyword},
        {key: 'Addr', value: this.keyword},
      ]
    }).subscribe(res => {

      this.channels = res.data;
      this.total = res.total;
    }, error => {
      console.log('error', error);
    }, () => {
      this.loading = false;
    });
  }

  loadFilters(): void {
    // this.as.get('distinct/copy/host').subscribe(res => {
    //   console.log('res', res);
    //   this.hosts = res.data.map(h => {
    //     return {
    //       text: h.host,
    //       value: h.host
    //     };
    //   });
    // }, error => {
    //   console.log('error', error);
    // });
  }

  edit(id?): void {
    const modal = this.ms.create({
      nzTitle: id ? '编辑通道' : '创建通道',
      nzContent: TunnelEditComponent,
      nzFooter: null,
      nzMaskClosable: false,
      // nzViewContainerRef: this.viewContainerRef,
      nzComponentParams: {id},
    });
    // insert/update after close
    modal.afterClose.subscribe(data => {
      if (!data) {
        return;
      }
      if (id) {
        this.channels.forEach((c: any, i, a: any[]) => {
            if (c.id === data.id) {
              a[i] = data;
            }
          }
        );
      } else {
        this.channels.unshift(data);
      }
    });
  }

  onTableQuery(params: NzTableQueryParams): void {
    const {pageSize, pageIndex, sort, filter} = params;
    this.pageSize = pageSize;
    this.pageIndex = pageIndex;
    const currentSort = sort.find(item => item.value !== null);
    this.sortField = (currentSort && currentSort.key) || null;
    this.sortOrder = (currentSort && currentSort.value) || null;
    this.filters = filter;
    this.load();
  }

  search(): void {
    this.pageIndex = 1;
    this.load();
  }
}
